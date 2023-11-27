package db

import (
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Connect() (*mongo.Client, error) {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	username := os.Getenv("MONGODB_USERNAME")
	password := os.Getenv("MONGODB_PASSWORD")
	uri := os.Getenv("MONGODB_URI")
	path := fmt.Sprintf(uri, username, password)
	opts := options.Client().ApplyURI(path).SetServerAPIOptions(serverAPI)
	return mongo.Connect(context.TODO(), opts)
}
