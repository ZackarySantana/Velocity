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
	"github.com/zackarysantana/velocity/internal/service/kafka"
	mongodomain "github.com/zackarysantana/velocity/internal/service/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	if os.Getenv("DEV_MODE") != "true" {
		logger.Debug("Loading env file")
		if err := godotenv.Load("env/.env.prod"); err != nil {
			log.Fatal("Error loading .env file", err)
		}
		logger.Debug("Loaded env file")
	}

	logger.Debug("Connecting to MongoDB")
	client, err := mongodomain.NewMongoClientFromEnv()
	if err != nil {
		panic(err)
	}
	repository := mongodomain.NewMongoRepositoryManager(client, os.Getenv("MONGODB_DATABASE"))
	logger.Debug("Connected to MongoDB")

	logger.Debug("Connecting to Kafka")
	pq, err := kafka.NewKafkaQueue(kafka.NewKafkaQueueOptionsFromEnv(os.Getenv("KAFKA_GROUP_ID_API")))
	defer pq.Close()
	if err != nil {
		panic(err)
	}
	logger.Debug("Connected to Kafka")

	serviceImpl := domain.NewService(repository, pq, mongodomain.NewMongoIdCreator(), logger)

	// delete
	// TODO: test of priority queue via mongodb.
	pqt := mongodomain.NewMongoPriorityQueue[string](client, mongodomain.NewMongoIdCreator(), os.Getenv("MONGODB_DATABASE"))
	// err = pqt.Push(context.TODO(), "test_queue", service.PriorityQueueItem[string]{Priority: 1, Payload: "testing this thing"}, service.PriorityQueueItem[string]{Priority: 2, Payload: "testing this thing 2"}, service.PriorityQueueItem[string]{Priority: 3, Payload: "testing this thing 3"})
	// fmt.Println(err)

	item, err := pqt.Pop(context.TODO(), "test_queue")
	fmt.Println(item, err)
	// delete

	mux := api.New[primitive.ObjectID](repository, serviceImpl, mongodomain.NewMongoIdCreator(), logger)
	logger.Info("Starting server", "addr", ":8080")
	http.ListenAndServe(":8080", mux)
}
