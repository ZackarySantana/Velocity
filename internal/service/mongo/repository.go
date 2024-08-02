package mongo

import (
	"context"

	"github.com/sysulq/dataloader-go"
	"github.com/zackarysantana/velocity/internal/service"
	"github.com/zackarysantana/velocity/src/entities/image"
	"github.com/zackarysantana/velocity/src/entities/job"
	"github.com/zackarysantana/velocity/src/entities/routine"
	"github.com/zackarysantana/velocity/src/entities/test"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	dbName = "velocity-beta"

	routineCollection = "routines"
	jobCollection     = "jobs"
	imageCollection   = "images"
	testCollection    = "tests"
)

func NewMongoRepository(db *mongo.Client) *service.Repository {
	return &service.Repository{
		Routine:     createTypeRepository[routine.Routine](db, dbName, routineCollection),
		PutRoutines: createPutType[routine.Routine](db, dbName, routineCollection),

		Job:     createTypeRepository[job.Job](db, dbName, jobCollection),
		PutJobs: createPutType[job.Job](db, dbName, jobCollection),

		Image:     createTypeRepository[image.Image](db, dbName, imageCollection),
		PutImages: createPutType[image.Image](db, dbName, imageCollection),

		Test:     createTypeRepository[test.Test](db, dbName, testCollection),
		PutTests: createPutType[test.Test](db, dbName, testCollection),

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

func createTypeRepository[T any](db *mongo.Client, database, collection string) dataloader.Interface[string, *T] {
	return dataloader.New(
		func(ctx context.Context, keys []string) []dataloader.Result[*T] {
			cur, err := db.Database(database).Collection(collection).Find(ctx, bson.M{
				"_id": bson.M{"$in": keys},
			})
			if err != nil {
				return []dataloader.Result[*T]{dataloader.Wrap[*T](nil, err)}
			}

			results := make([]dataloader.Result[*T], len(keys))
			for cur.Next(ctx) {
				var r *T
				if err := cur.Decode(&r); err != nil {
					return []dataloader.Result[*T]{dataloader.Wrap[*T](nil, err)}
				}
				results = append(results, dataloader.Wrap(r, nil))
			}

			return results
		},
	)
}

func createPutType[T any](db *mongo.Client, database, collection string) func(context.Context, []*T) error {
	return func(ctx context.Context, t []*T) error {
		items := make([]interface{}, len(t))
		for i, v := range t {
			items[i] = v
		}
		_, err := db.Database(database).Collection(collection).InsertMany(ctx, items)
		return err
	}
}
