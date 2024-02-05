package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/zackarysantana/velocity/internal/api"
	"github.com/zackarysantana/velocity/internal/cli/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	db, err := getEnv("MONGODB_DATABASE")
	if err != nil {
		log.Fatal(err)
	}

	username, err := getEnv("MONGODB_USERNAME")
	if err != nil {
		log.Fatal(err)
	}

	password, err := getEnv("MONGODB_PASSWORD")
	if err != nil {
		log.Fatal(err)
	}

	uri, err := getEnv("MONGODB_URI")
	if err != nil {
		log.Fatal(err)
	}
	path := fmt.Sprintf(uri, username, password)
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(path).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		log.Fatal(err)
	}

	l := logger.NewLiveLogger()
	l.SubscribeError(os.Stdout)

	engine := api.CreateApi(l, client)
	engine.AddAgentRoutes(db, "agents")
	engine.Run(":8080")
}

func getEnv(name string) (string, error) {
	value, exists := os.LookupEnv(name)
	if !exists {
		return "", fmt.Errorf("%s not set", name)
	}
	return value, nil
}
