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
	"github.com/zackarysantana/velocity/internal/service"
	"github.com/zackarysantana/velocity/internal/service/domain"
	"github.com/zackarysantana/velocity/internal/service/kafka"
	mongodomain "github.com/zackarysantana/velocity/internal/service/mongo"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file", err)
	}

	repository, err := loadMongoDB()
	if err != nil {
		panic(err)
	}

	pq, err := kafka.NewKafkaQueue(os.Getenv("KAFKA_USERNAME"), os.Getenv("KAFKA_PASSWORD"), os.Getenv("KAFKA_BROKER"), "agent")
	defer pq.Close()
	if err != nil {
		panic(err)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	service := domain.NewService(repository, pq, logger)

	mux := api.New(repository, service, mongodomain.NewMongoIdCreator(), logger)
	slog.Info("Starting server", "addr", ":8080")
	http.ListenAndServe(":8080", mux)
}

func loadMongoDB() (*service.Repository, error) {
	uri := fmt.Sprintf(os.Getenv("MONGODB_URI"), os.Getenv("MONGODB_USERNAME"), os.Getenv("MONGODB_PASSWORD"))
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	return mongodomain.NewMongoRepository(client, os.Getenv("MONGODB_DATABASE")), nil
}
