package mongo

import (
	"context"
	"fmt"
	"os"

	"github.com/zackarysantana/velocity/internal/service"
	"github.com/zackarysantana/velocity/src/entities/image"
	"github.com/zackarysantana/velocity/src/entities/job"
	"github.com/zackarysantana/velocity/src/entities/routine"
	"github.com/zackarysantana/velocity/src/entities/test"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	routineCollection = "routines"
	jobCollection     = "jobs"
	imageCollection   = "images"
	testCollection    = "tests"
)

func NewMongoClientFromEnv() (*mongo.Client, error) {
	uri := fmt.Sprintf(os.Getenv("MONGODB_URI"), os.Getenv("MONGODB_USERNAME"), os.Getenv("MONGODB_PASSWORD"))
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	if err := client.Ping(context.Background(), nil); err != nil {
		return nil, err
	}
	return client, nil
}

func NewMongoRepositoryManager(db *mongo.Client, dbName string) *service.RepositoryManager[primitive.ObjectID] {
	return &service.RepositoryManager[primitive.ObjectID]{
		Routine: &service.RoutineRepository[primitive.ObjectID]{
			Load: createLoad[routine.Routine[primitive.ObjectID]](db, dbName, routineCollection),
			Put:  createPutForType[routine.Routine[primitive.ObjectID]](db, dbName, routineCollection),
		},
		Job: &service.JobRepository[primitive.ObjectID]{
			Load: createLoad[job.Job[primitive.ObjectID]](db, dbName, jobCollection),
			Put:  createPutForType[job.Job[primitive.ObjectID]](db, dbName, jobCollection),
		},
		Image: &service.ImageRepository[primitive.ObjectID]{
			Load: createLoad[image.Image[primitive.ObjectID]](db, dbName, imageCollection),
			Put:  createPutForType[image.Image[primitive.ObjectID]](db, dbName, imageCollection),
		},
		Test: &service.TestRepository[primitive.ObjectID]{
			Load: createLoad[test.Test[primitive.ObjectID]](db, dbName, testCollection),
			Put:  createPutForType[test.Test[primitive.ObjectID]](db, dbName, testCollection),
		},
		WithTransaction: func(ctx context.Context, fn func(context.Context) error) error {
			session, err := db.StartSession()
			if err != nil {
				return err
			}
			defer session.EndSession(ctx)
			_, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
				return nil, fn(sessCtx)
			})
			return err
		},
	}
}

func createLoad[T any](db *mongo.Client, database, collection string) func(context.Context, []primitive.ObjectID) ([]*T, error) {
	return func(ctx context.Context, keys []primitive.ObjectID) ([]*T, error) {
		cur, err := db.Database(database).Collection(collection).Find(ctx, bson.M{
			"_id": bson.M{"$in": keys},
		})
		if err != nil {
			return nil, err
		}

		results := make([]*T, cur.RemainingBatchLength())
		i := 0
		for cur.Next(ctx) {
			var r T
			if err := cur.Decode(&r); err != nil {
				return nil, err
			}
			results[i] = &r
			i++
		}

		return results, nil
	}
}

func createPutForType[T any](db *mongo.Client, database, collection string) func(context.Context, []*T) error {
	return func(ctx context.Context, t []*T) error {
		items := make([]interface{}, len(t))
		for i, v := range t {
			items[i] = v
		}
		_, err := db.Database(database).Collection(collection).InsertMany(ctx, items)
		return err
	}
}
