package mongo

import (
	"context"
	"fmt"
	"os"

	"github.com/sysulq/dataloader-go"
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

func NewMongoRepository(db *mongo.Client, dbName string) *service.Repository {
	return &service.Repository{
		Routine: service.RoutineRepository{
			Interface: createTypeRepository[routine.Routine](db, dbName, routineCollection),
			Put:       createPutType[routine.Routine](db, dbName, routineCollection),
		},
		Job: service.JobRepository{
			Interface: createTypeRepository[job.Job](db, dbName, jobCollection),
			Put:       createPutType[job.Job](db, dbName, jobCollection),
		},
		Image: service.ImageRepository{
			Interface: createTypeRepository[image.Image](db, dbName, imageCollection),
			Put:       createPutType[image.Image](db, dbName, imageCollection),
		},
		Test: service.TestRepository{
			Interface: createTypeRepository[test.Test](db, dbName, testCollection),
			Put:       createPutType[test.Test](db, dbName, testCollection),
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
