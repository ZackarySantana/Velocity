package domain

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

func NewRepository(db *mongo.Client) (*service.Repository, error) {
	return &service.Repository{
		Routine:    createTypeRepository[routine.Routine](db, "velocity", "routines"),
		PutRoutine: createPutType[routine.Routine](db, "velocity", "routines"),

		Job:    createTypeRepository[job.Job](db, "velocity", "jobs"),
		PutJob: createPutType[job.Job](db, "velocity", "jobs"),

		Image:    createTypeRepository[image.Image](db, "velocity", "images"),
		PutImage: createPutType[image.Image](db, "velocity", "images"),

		Test:    createTypeRepository[test.Test](db, "velocity", "tests"),
		PutTest: createPutType[test.Test](db, "velocity", "tests"),
	}, nil
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

func createPutType[T any](db *mongo.Client, database, collection string) func(context.Context, *T) error {
	return func(ctx context.Context, t *T) error {
		_, err := db.Database(database).Collection(collection).InsertOne(ctx, t)
		return err
	}
}
