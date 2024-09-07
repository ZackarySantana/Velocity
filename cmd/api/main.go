package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/joho/godotenv"
	"github.com/zackarysantana/velocity/cmd/api/internal"
	"github.com/zackarysantana/velocity/internal/api"
	"github.com/zackarysantana/velocity/internal/otel"
	"github.com/zackarysantana/velocity/internal/service"
	"github.com/zackarysantana/velocity/internal/service/domain"
)

func main() {
	// Handle SIGINT (CTRL+C) gracefully.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	if os.Getenv("DEV_MODE") != "true" {
		logger.Debug("Loading env file")
		if err := godotenv.Load("env/.env.prod"); err != nil {
			logger.Debug("Error loading .env file", "error", err)
		}
		logger.Debug("Loaded env file")
	}

	idCreator := internal.GetIDCreator[any](logger)

	logger.Debug("Connecting to repository manager...")
	repository := internal.GetRepositoryManager(logger, idCreator)
	logger.Debug("Connected to repository manager")

	logger.Debug("Connecting to process queue...")
	pq := internal.GetProcessQueue(logger)
	defer pq.Close()
	logger.Debug("Connected to process queue")

	serviceImpl := domain.NewService(repository, pq, idCreator, logger)

	logger.Debug("Connecting to priority queue...")
	pqt := internal.GetPriorityQueue[any, string](logger)
	logger.Debug("Connected to priority queue")

	pqt.Push(ctx, "test_queue", service.PriorityQueueItem[string]{Priority: 1, Payload: "1"})
	fmt.Println(pqt.Pop(ctx, "test_queue"))

	shutdown, err := otel.Setup(ctx)
	defer shutdown(ctx)
	if err != nil {
		panic(err)
	}

	mux := api.New(repository, serviceImpl, idCreator, logger)
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
