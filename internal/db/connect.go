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

	db string
}

func (c *Connection) ApplyIndexes(ctx context.Context) error {
	i := []func(context.Context) error{
		c.ApplyUserIndexes,
		c.ApplyPermissionIndexes,
		c.ApplyProjectIndexes,
		c.ApplyInstanceIndexes,
		c.ApplyJobIndexes,
	}

	errs := []error{}
	for _, f := range i {
		if e := f(ctx); e != nil {
			errs = append(errs, e)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors applying indexes: %v", errs)
	}

	return nil
}

func Connect(ctx *context.Context) (*Connection, error) {
	if ctx == nil {
		defaultContext := context.TODO()
		ctx = &defaultContext
	}

	db, err := getEnv("MONGODB_DATABASE")
	if err != nil {
		return nil, err
	}

	username, err := getEnv("MONGODB_USERNAME")
	if err != nil {
		return nil, err
	}

	password, err := getEnv("MONGODB_PASSWORD")
	if err != nil {
		return nil, err
	}

	uri, err := getEnv("MONGODB_URI")
	if err != nil {
		return nil, err
	}
	path := fmt.Sprintf(uri, username, password)
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(path).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(*ctx, opts)
	if err != nil {
		return nil, err
	}

	return &Connection{client, db}, err
}

func getEnv(name string) (string, error) {
	value, exists := os.LookupEnv(name)
	if !exists {
		return "", fmt.Errorf("%s not set", name)
	}
	return value, nil
}

func (c *Connection) col(collection string) *mongo.Collection {
	return c.Database(c.db).Collection(collection)
}
