package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/zackarysantana/velocity/internal/api"
	"github.com/zackarysantana/velocity/internal/service/domain"
	mongodomain "github.com/zackarysantana/velocity/internal/service/mongo"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file", err)
	}

	uri := fmt.Sprintf(os.Getenv("MONGODB_URI"), os.Getenv("MONGODB_USERNAME"), os.Getenv("MONGODB_PASSWORD"))
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))

	if err != nil {
		log.Fatal("Error connecting to MongoDB", err)
	}
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal("Error pinging MongoDB", err)
	}

	dbName := os.Getenv("MONGODB_DATABASE")
	mux := api.New(domain.NewService(mongodomain.NewMongoRepository(client, dbName)), mongodomain.NewMongoIdCreator())
	slog.Info("Starting server", "addr", ":8080")
	http.ListenAndServe(":8080", mux)
}
