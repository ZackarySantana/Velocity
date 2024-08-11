package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/zackarysantana/velocity/internal/api"
	"github.com/zackarysantana/velocity/internal/service/domain"
	"github.com/zackarysantana/velocity/internal/service/kafka"
	mongodomain "github.com/zackarysantana/velocity/internal/service/mongo"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	logger.Debug("Loading env file")
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file", err)
	}
	logger.Debug("Loaded env file")

	logger.Debug("Connecting to MongoDB")
	client, err := mongodomain.NewMongoClientFromEnv()
	if err != nil {
		panic(err)
	}
	repository := mongodomain.NewMongoRepository(client, os.Getenv("MONGODB_DATABASE"))
	logger.Debug("Connected to MongoDB")

	logger.Debug("Connecting to Kafka")
	pq, err := kafka.NewKafkaQueue(kafka.NewKafkaQueueOptionsFromEnv())
	defer pq.Close()
	if err != nil {
		panic(err)
	}
	logger.Debug("Connected to Kafka")

	service := domain.NewService(repository, pq, logger)

	mux := api.New(repository, service, mongodomain.NewMongoIdCreator(), logger)
	logger.Info("Starting server", "addr", ":8080")
	http.ListenAndServe(":8080", mux)
}
