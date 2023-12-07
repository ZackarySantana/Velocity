package db

import (
	"context"
	"fmt"
	"time"

	"github.com/zackarysantana/velocity/internal/jobs/jobtypes"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Job struct {
	Id primitive.ObjectID `bson:"_id,omitempty" json:"id"`

	// Run logic
	Name          string   `bson:"name,omitempty" json:"name"`
	Image         string   `bson:"image,omitempty" json:"image"`
	Command       string   `bson:"command,omitempty" json:"command"`
	SetupCommands []string `bson:"setup_commands,omitempty" json:"setup_commands"`

	// Status logic
	Status jobtypes.JobStatus `bson:"status,omitempty" json:"status"`
	Logs   string             `bson:"logs,omitempty" json:"logs"`
	Error  string             `bson:"error,omitempty" json:"error"`

	// For tasks attached to an instance only
	InstanceId *primitive.ObjectID `bson:"instance_id,omitempty" json:"instance_id,omit_empty"`

	// Metadata
	CreatedAt time.Time `bson:"created_at,omitempty" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at,omitempty" json:"updated_at"`
}

func (j *Job) validate() error {
	if !j.Id.IsZero() {
		return nil
	}
	if j.Name == "" {
		return fmt.Errorf("missing name")
	}
	if j.Image == "" {
		return fmt.Errorf("missing image on %s", j.Name)
	}
	if j.Command == "" {
		return fmt.Errorf("missing command %s", j.Name)
	}
	if j.Status == "" {
		return fmt.Errorf("missing status %s", j.Name)
	}
	return nil
}

func (c *Connection) GetJob(ctx context.Context, query interface{}, opts ...*options.FindOneOptions) (*Job, error) {
	var job Job
	return &job, c.col("jobs").FindOne(ctx, query, opts...).Decode(&job)
}

func (c *Connection) GetJobs(ctx context.Context, query interface{}, opts ...*options.FindOptions) ([]*Job, error) {
	var jobs []*Job
	cur, err := c.col("jobs").Find(ctx, query, opts...)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var job Job
		if err := cur.Decode(&job); err != nil {
			return nil, err
		}
		jobs = append(jobs, &job)
	}
	return jobs, nil
}

func (c *Connection) GetNQueuedJobs(ctx context.Context, n int64) ([]*Job, error) {
	query := bson.M{"status": jobtypes.JobStatusQueued}
	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetLimit(n)
	return c.GetJobs(ctx, query, opts)
}

func (c *Connection) DequeueNJobs(ctx context.Context, n int64) ([]*Job, error) {
	jobs := []*Job{}
	err := c.UseSessionWithOptions(ctx, nil, func(ctx mongo.SessionContext) error {
		var err error
		if err = ctx.StartTransaction(); err != nil {
			return err
		}

		jobs, err = c.GetNQueuedJobs(ctx, n)
		if err != nil {
			_ = ctx.AbortTransaction(context.Background())
			return err
		}

		ids := []primitive.ObjectID{}
		for _, job := range jobs {
			ids = append(ids, job.Id)
		}

		update := bson.M{"$set": bson.M{"status": jobtypes.JobStatusActive}}
		_, err = c.UpdateJobsByIds(ctx, ids, nil, update)
		if err != nil {
			_ = ctx.AbortTransaction(context.Background())
			return err
		}

		return ctx.CommitTransaction(context.Background())
	})
	if err != nil {
		return nil, err
	}

	return jobs, nil
}

func (c *Connection) InsertJob(ctx context.Context, job *Job) (*Job, error) {
	if err := job.validate(); err != nil {
		return nil, err
	}
	job.CreatedAt = time.Now()
	job.UpdatedAt = time.Now()
	r, err := c.col("jobs").InsertOne(ctx, job)
	if err != nil {
		return nil, err
	}
	job.Id = r.InsertedID.(primitive.ObjectID)

	return job, nil
}

func (c *Connection) UpdateJob(ctx context.Context, job *Job) (*Job, error) {
	if err := job.validate(); err != nil {
		return nil, err
	}
	job.UpdatedAt = time.Now()
	_, err := c.col("jobs").UpdateOne(ctx, bson.M{"_id": job.Id}, bson.M{"$set": job})
	if err != nil {
		return nil, err
	}
	return job, nil
}

func (c *Connection) InsertJobs(ctx context.Context, jobs []*Job) ([]*Job, error) {
	for _, job := range jobs {
		if err := job.validate(); err != nil {
			return nil, err
		}
		job.CreatedAt = time.Now()
		job.UpdatedAt = time.Now()
	}
	var docs []interface{}
	for _, job := range jobs {
		docs = append(docs, job)
	}
	r, err := c.col("jobs").InsertMany(ctx, docs)
	if err != nil {
		return nil, err
	}
	for i, id := range r.InsertedIDs {
		jobs[i].Id = id.(primitive.ObjectID)
	}
	return jobs, nil
}

func (c *Connection) UpdateJobsByIds(ctx context.Context, ids []primitive.ObjectID, filter, update bson.M) (int, error) {
	if filter == nil {
		filter = bson.M{}
	}
	filter["_id"] = bson.M{"$in": ids}
	if update == nil {
		update = bson.M{}
	}
	if _, ok := update["$set"]; !ok {
		update["$set"] = bson.M{}
	}
	update["$set"].(bson.M)["updated_at"] = time.Now()
	r, err := c.col("jobs").UpdateMany(ctx, filter, update)
	if err != nil {
		return -1, err
	}
	return int(r.MatchedCount), nil
}
