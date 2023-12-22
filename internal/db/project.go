package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Project struct {
	Id primitive.ObjectID `bson:"_id,omitempty" json:"id"`

	Name string `bson:"name,omitempty" json:"name"`

	Git struct {
		Owner      string `bson:"owner,omitempty" json:"owner"`
		Repository string `bson:"repository,omitempty" json:"repository"`
	} `bson:"git,omitempty" json:"git"`
}

func (c *Connection) ApplyProjectIndexes(ctx context.Context) error {
	_, err := c.col("projects").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.M{"name": 1},
	})
	return err
}

func (c *Connection) GetProject(ctx context.Context, query interface{}) (*Project, error) {
	var project Project
	return &project, c.col("projects").FindOne(ctx, query).Decode(&project)
}

func (c *Connection) GetProjectById(ctx context.Context, id primitive.ObjectID) (*Project, error) {
	return c.GetProject(ctx, bson.M{"_id": id})
}

// By name
func (c *Connection) GetProjectByName(ctx context.Context, name string) (*Project, error) {
	return c.GetProject(ctx, bson.M{"name": name})
}

func (c *Connection) InsertProject(ctx context.Context, project *Project) (*Project, error) {
	r, err := c.col("projects").InsertOne(ctx, project)
	if err != nil {
		return nil, err
	}
	project.Id = r.InsertedID.(primitive.ObjectID)

	return project, nil
}
