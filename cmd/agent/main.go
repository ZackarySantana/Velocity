package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/zackarysantana/velocity/internal/agent"
	"github.com/zackarysantana/velocity/internal/service/kafka"
	"github.com/zackarysantana/velocity/src/velocity"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	logger.Debug("Loading env file")
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file", err)
	}
	logger.Debug("Loaded env file")

	logger.Debug("Connecting to Kafka")
	pq, err := kafka.NewKafkaQueue(kafka.NewKafkaQueueOptionsFromEnv())
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
	agent := agent.New(pq, velocity)
	err = agent.Start(context.Background())
	if err != nil {
		panic(err)
	}
}
