package db

import (
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Connection struct {
	*mongo.Client
}

func Connect(ctx *context.Context) (*Connection, error) {
	if ctx == nil {
		defaultContext := context.TODO()
		ctx = &defaultContext
	}
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	username := os.Getenv("MONGODB_USERNAME")
	password := os.Getenv("MONGODB_PASSWORD")
	uri := os.Getenv("MONGODB_URI")
	path := fmt.Sprintf(uri, username, password)
	opts := options.Client().ApplyURI(path).SetServerAPIOptions(serverAPI)
	client, err := mongo.Connect(*ctx, opts)
	return &Connection{client}, err
}
