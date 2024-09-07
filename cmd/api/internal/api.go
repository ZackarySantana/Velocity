package internal

import (
	"context"
	"log/slog"
	"os"

	"github.com/zackarysantana/velocity/internal/service"
	"github.com/zackarysantana/velocity/internal/service/kafka"
	"github.com/zackarysantana/velocity/internal/service/mock"
	mongodomain "github.com/zackarysantana/velocity/internal/service/mongo"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetIDCreator[T any](logger *slog.Logger) service.IDCreator[T] {
	useMock := os.Getenv("MOCK_ID_CREATOR")
	if useMock == "true" {
		logger.Debug("Using mock ID creator")
		return mock.NewMockIDCreator[T]()
	}
	useMongo := os.Getenv("MONGO_ID_CREATOR")
	if useMongo == "true" {
		logger.Debug("Using mongo ID creator")
		return mongodomain.NewObjectIDCreator[T]().(service.IDCreator[T])
	}

	panic("No ID creator set")
}

func GetRepositoryManager[T comparable](logger *slog.Logger, idCreator service.IDCreator[T]) service.RepositoryManager[T] {
	useMock := os.Getenv("MOCK_REPOSITORY_MANAGER")
	if useMock == "true" {
		logger.Debug("Using mock repository manager")
		return mock.NewMockRepositoryManager[T](GetIDCreator[T](logger))
	}
	useMongo := os.Getenv("MONGO_REPOSITORY_MANAGER")
	if useMongo == "true" {
		logger.Debug("Using mongo repository manager")
		client, err := mongo.Connect(context.Background(), mongodomain.URIFromEnv())
		if err != nil {
			panic(err)
		}
		if err := client.Ping(context.Background(), nil); err != nil {
			panic(err)
		}
		return mongodomain.NewMongoRepositoryManager[T](client, os.Getenv("MONGODB_DATABASE"))
	}

	panic("No repository manager set")
}

func GetProcessQueue(logger *slog.Logger) service.ProcessQueue {
	useKafka := os.Getenv("KAFKA_PROCESS_QUEUE")
	if useKafka == "true" {
		logger.Debug("Using kafka process queue")
		pq, err := kafka.NewKafkaQueue(kafka.NewKafkaQueueOptionsFromEnv(os.Getenv("KAFKA_GROUP_ID_API")))
		if err != nil {
			panic(err)
		}
		return pq
	}

	panic("No process queue set")
}

func GetPriorityQueue[ID any, Payload any](logger *slog.Logger) service.PriorityQueue[ID, Payload] {
	useMongo := os.Getenv("MONGO_PRIORITY_QUEUE")
	if useMongo == "true" {
		logger.Debug("Using mongo priority queue")
		client, err := mongo.Connect(context.Background(), mongodomain.URIFromEnv())
		if err != nil {
			panic(err)
		}
		if err := client.Ping(context.Background(), nil); err != nil {
			panic(err)
		}
		idCreator := GetIDCreator[ID](logger)
		return mongodomain.NewMongoPriorityQueue[ID, Payload](client, idCreator, os.Getenv("MONGODB_DATABASE"))
	}

	panic("No priority queue set")
}
