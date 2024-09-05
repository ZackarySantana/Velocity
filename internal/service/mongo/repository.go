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
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	routineCollection = "routines"
	jobCollection     = "jobs"
	imageCollection   = "images"
	testCollection    = "tests"
)

func URIFromEnv() *options.ClientOptions {
	uri := fmt.Sprintf(os.Getenv("MONGODB_URI"), os.Getenv("MONGODB_USERNAME"), os.Getenv("MONGODB_PASSWORD"))
	return options.Client().ApplyURI(uri)
}

// NewMongoRepositoryManager creates a new RepositoryManager with the provided mongo client and database name.
// The type provided is used as the Key type for the documents.
func NewMongoRepositoryManager[T any](db *mongo.Client, dbName string) *service.RepositoryManager[T] {
	return &service.RepositoryManager[T]{
		Routine: &service.RoutineRepository[T]{
			Load: createLoad[routine.Routine[T], T](db, dbName, routineCollection),
			Put:  createPutForType[routine.Routine[T]](db, dbName, routineCollection),
		},
		Job: &service.JobRepository[T]{
			Load: createLoad[job.Job[T], T](db, dbName, jobCollection),
			Put:  createPutForType[job.Job[T]](db, dbName, jobCollection),
		},
		Image: &service.ImageRepository[T]{
			Load: createLoad[image.Image[T], T](db, dbName, imageCollection),
			Put:  createPutForType[image.Image[T]](db, dbName, imageCollection),
		},
		Test: &service.TestRepository[T]{
			Load: createLoad[test.Test[T], T](db, dbName, testCollection),
			Put:  createPutForType[test.Test[T]](db, dbName, testCollection),
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

func createLoad[T any, V any](db *mongo.Client, database, collection string) func(context.Context, []V) ([]*T, error) {
	return func(ctx context.Context, keys []V) ([]*T, error) {
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
