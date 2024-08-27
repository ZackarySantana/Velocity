package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/joho/godotenv"
	"github.com/zackarysantana/velocity/internal/api"
	"github.com/zackarysantana/velocity/internal/otel"
	"github.com/zackarysantana/velocity/internal/service/domain"
	"github.com/zackarysantana/velocity/internal/service/kafka"
	mongodomain "github.com/zackarysantana/velocity/internal/service/mongo"
)

func main() {
	// Handle SIGINT (CTRL+C) gracefully.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

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

	item, err := pqt.Pop(ctx, "test_queue")
	fmt.Println(item, err)
	// delete

	shutdown, err := otel.Setup(ctx)
	defer shutdown(ctx)
	if err != nil {
		panic(err)
	}

	mux := api.New(repository, serviceImpl, mongodomain.NewMongoIdCreator(), logger)
	logger.Info("Starting server", "addr", ":8080")
	srv := &http.Server{
		Addr:         ":8080",
		BaseContext:  func(_ net.Listener) context.Context { return ctx },
		ReadTimeout:  time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      mux,
	}
	srvErr := make(chan error, 1)
	go func() {
		srvErr <- srv.ListenAndServe()
	}()

	// Wait for interruption.
	select {
	case err = <-srvErr:
		// Error when starting HTTP server.
		return
	case <-ctx.Done():
		// Wait for first CTRL+C.
		// Stop receiving signal notifications as soon as possible.
		stop()
	}

	// When Shutdown is called, ListenAndServe immediately returns ErrServerClosed.
	err = srv.Shutdown(context.Background())
	return
}
