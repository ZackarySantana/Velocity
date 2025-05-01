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
	"github.com/zackarysantana/velocity/cmd/internal"
	"github.com/zackarysantana/velocity/internal/api"
	"github.com/zackarysantana/velocity/internal/otel"
	"github.com/zackarysantana/velocity/internal/service/domain"
)

func main() {
	// Handle SIGINT (CTRL+C) gracefully.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	if os.Getenv("DEV_MODE") == "true" && os.Getenv("DEV_SERVICES") != "true" {
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

	// ProcessQueue no longer in use.
	// logger.Debug("Connecting to process queue...")
	// pq := internal.GetProcessQueue(logger)
	// defer pq.Close()
	// logger.Debug("Connected to process queue")

	logger.Debug("Connecting to priority queue...")
	pq := internal.GetPriorityQueue[any, any](logger)
	logger.Debug("Connected to priority queue")

	serviceImpl := domain.NewService(idCreator, repository, pq, logger)

	shutdown, err := otel.Setup(ctx)
	defer shutdown(ctx)
	if err != nil {
		panic(err)
	}

	port := ":8080"
	if os.Getenv("PORT") != "" {
		port = ":" + os.Getenv("PORT")
	}

	mux := api.New(idCreator, repository, serviceImpl, pq, logger)
	logger.Info("Starting server", "addr", port)
	srv := &http.Server{
		Addr:         port,
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
		panic(err)
	case <-ctx.Done():
		// Wait for first CTRL+C.
		// Stop receiving signal notifications as soon as possible.
		fmt.Println("Shutting down server...")
		stop()
	}

	// When Shutdown is called, ListenAndServe immediately returns ErrServerClosed.
	err = srv.Shutdown(context.Background())
	if err != nil {
		panic(err)
	}
}
