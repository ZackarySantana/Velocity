package mongotest

import (
	"context"

	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateContainer(ctx context.Context) (*mongo.Client, func(context.Context) error, error) {
	mongodbContainer, err := mongodb.Run(ctx, "mongo:6")
	if err != nil {
		return nil, nil, err
	}
	cleanup := func(ctx context.Context) error {
		return mongodbContainer.Terminate(ctx)
	}
	endpoint, err := mongodbContainer.ConnectionString(ctx)
	if err != nil {
		cleanup(ctx)
		return nil, nil, err
	}
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(endpoint))
	if err != nil {
		cleanup(ctx)
		return nil, nil, err
	}
	return mongoClient, cleanup, nil
}
