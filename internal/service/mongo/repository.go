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

func NewMongoRepositoryManager[ID any](client *mongo.Client, database string) service.RepositoryManager[ID] {
	return service.NewRepositoryManager(
		newExampleMongoRepository[ID, routine.Routine[ID]](client, database, routineCollection),
		newExampleMongoRepository[ID, job.Job[ID]](client, database, jobCollection),
		newExampleMongoRepository[ID, image.Image[ID]](client, database, imageCollection),
		newExampleMongoRepository[ID, test.Test[ID]](client, database, testCollection),
		func(ctx context.Context, fn func(context.Context) error) error {
			session, err := client.StartSession()
			if err != nil {
				return err
			}
			defer session.EndSession(ctx)
			_, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
				return nil, fn(sessCtx)
			})
			return err
		},
	)
}

func newExampleMongoRepository[ID any, DataType any](client *mongo.Client, database string, collection string) service.TypeRepository[ID, DataType] {
	return &exampleMongoRepository[ID, DataType]{
		client:     client,
		database:   database,
		collection: collection,
	}
}

type exampleMongoRepository[ID any, DataType any] struct {
	client     *mongo.Client
	database   string
	collection string
}

func (e *exampleMongoRepository[ID, DataType]) Load(ctx context.Context, keys []ID) ([]*DataType, error) {
	cur, err := e.client.Database(e.database).Collection(e.collection).Find(ctx, bson.M{
		"_id": bson.M{"$in": keys},
	})
	if err != nil {
		return nil, err
	}

	results := make([]*DataType, cur.RemainingBatchLength())
	i := 0
	for cur.Next(ctx) {
		var r DataType
		if err := cur.Decode(&r); err != nil {
			return nil, err
		}
		results[i] = &r
		i++
	}

	return results, nil
}

func (e *exampleMongoRepository[ID, DataType]) Put(ctx context.Context, data []*DataType) ([]ID, error) {
	items := make([]interface{}, len(data))
	for i, v := range data {
		items[i] = v
	}
	insertedIDs, err := e.client.Database(e.database).Collection(e.collection).InsertMany(ctx, items)
	if err != nil {
		return nil, err
	}
	keys := make([]ID, len(insertedIDs.InsertedIDs))
	for i, id := range insertedIDs.InsertedIDs {
		keys[i] = id.(ID)
	}
	return keys, nil
}
