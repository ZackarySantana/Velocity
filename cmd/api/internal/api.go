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
		return mock.NewIDCreator[T]()
	}
	useMongo := os.Getenv("MONGO_ID_CREATOR")
	if useMongo == "true" {
		logger.Debug("Using mongo ID creator")
		return mongodomain.NewIDCreator[T]().(service.IDCreator[T])
	}

	panic("No ID creator set")
}

func GetRepositoryManager[T comparable](logger *slog.Logger, idCreator service.IDCreator[T]) service.RepositoryManager[T] {
	useMock := os.Getenv("MOCK_REPOSITORY_MANAGER")
	if useMock == "true" {
		logger.Debug("Using mock repository manager")
		return mock.NewRepositoryManager[T](GetIDCreator[T](logger))
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
		return mongodomain.NewRepositoryManager[T](client, os.Getenv("MONGODB_DATABASE"))
	}

	panic("No repository manager set")
}

func GetProcessQueue(logger *slog.Logger) service.ProcessQueue {
	useKafka := os.Getenv("KAFKA_PROCESS_QUEUE")
	if useKafka == "true" {
		logger.Debug("Using kafka process queue")
		pq, err := kafka.NewProcessQueue(kafka.NewProcessQueueConfigFromEnv(os.Getenv("KAFKA_GROUP_ID_API")))
		if err != nil {
			panic(err)
		}
		return pq
	}

	panic("No process queue set")
}

func GetPriorityQueue[ID comparable, Payload any](logger *slog.Logger) service.PriorityQueue[ID, Payload] {
	useMock := os.Getenv("MOCK_PRIORITY_QUEUE")
	if useMock == "true" {
		logger.Debug("Using mock priority queue")
		return mock.NewPriorityQueue[ID, Payload](GetIDCreator[ID](logger))
	}
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
		return mongodomain.NewPriorityQueue[ID, Payload](client, idCreator, os.Getenv("MONGODB_DATABASE"))
	}

	panic("No priority queue set")
}