package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"github.com/joho/godotenv"
	"github.com/zackarysantana/velocity/internal/agent"
	"github.com/zackarysantana/velocity/internal/service/kafka"
	"github.com/zackarysantana/velocity/src/velocity"
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

	logger.Debug("Connecting to Kafka")
	pq, err := kafka.NewProcessQueue(kafka.NewProcessQueueConfigFromEnv(os.Getenv("KAFKA_GROUP_ID_AGENT")))
	defer pq.Close()
	if err != nil {
		panic(err)
	}
	logger.Debug("Connected to Kafka")

	velocity := velocity.NewAgent(os.Getenv("VELOCITY_URL"))

	// If this is in dev mode, we wait a little because
	// both services usually get restarted at the same time.
	if os.Getenv("DEV_MODE") == "true" {
		time.Sleep(2 * time.Second)
	}

	logger.Debug("Starting agent")
	agent := agent.New(pq, velocity, logger)
	err = agent.Start(ctx)
	if err != nil {
		panic(err)
	}
}
