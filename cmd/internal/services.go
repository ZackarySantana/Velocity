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
	idCreator := os.Getenv("ID_CREATOR")
	if idCreator == "mock" {
		logger.Debug("Using mock ID creator")
		return mock.NewIDCreator[T]()
	}
	if idCreator == "mongo" {
		logger.Debug("Using mongo ID creator")
		return mongodomain.NewIDCreator[T]().(service.IDCreator[T])
	}

	panic("No ID creator set")
}

func GetRepositoryManager[T comparable](logger *slog.Logger, idCreator service.IDCreator[T]) service.RepositoryManager[T] {
	repositoryManager := os.Getenv("REPOSITORY_MANAGER")
	if repositoryManager == "mock" {
		logger.Debug("Using mock repository manager")
		return mock.NewRepositoryManager[T](GetIDCreator[T](logger))
	}
	if repositoryManager == "mongo" {
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
	processQueue := os.Getenv("PROCESS_QUEUE")
	if processQueue == "kafka" {
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
	priorityQueue := os.Getenv("PRIORITY_QUEUE")
	if priorityQueue == "mock" {
		logger.Debug("Using mock priority queue")
		return mock.NewPriorityQueue[ID, Payload](GetIDCreator[ID](logger))
	}
	if priorityQueue == "mongo" {
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
