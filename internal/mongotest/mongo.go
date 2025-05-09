package mongotest

import (
	"context"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CreateContainer creates a MongoDB container. It creates a replica set
// to support transactions.
func CreateContainer(ctx context.Context) (*mongo.Client, func(context.Context) error, error) {
	var reuse testcontainers.CustomizeRequestOption = func(req *testcontainers.GenericContainerRequest) error {
		req.Name = "mongo-test"
		req.Reuse = true
		return nil
	}

	mongodbContainer, err := mongodb.Run(ctx, "mongo:6", mongodb.WithReplicaSet("rs0"),
		reuse,
	)
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
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(endpoint+"/?replicaSet=&directConnection=true"))
	if err != nil {
		cleanup(ctx)
		return nil, nil, err
	}

	return mongoClient, cleanup, nil
}
