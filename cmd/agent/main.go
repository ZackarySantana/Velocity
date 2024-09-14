package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"github.com/joho/godotenv"
	"github.com/zackarysantana/velocity/cmd/internal"
	"github.com/zackarysantana/velocity/internal/agent"
	"github.com/zackarysantana/velocity/src/velocity"
)

var (
	agentGroupID = "velocity-agent"
)

func main() {
	// Handle SIGINT (CTRL+C) gracefully.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	if os.Getenv("DEV_MODE") == "true" && os.Getenv("DEV_SERVICES") != "true" {
		logger.Debug("Loading env file")
		if err := godotenv.Load("env/.env.prod"); err != nil {
			log.Fatal("Error loading .env file", err)
		}
		logger.Debug("Loaded env file")
	}

	logger.Debug("Connecting to process queue...")
	pq := internal.GetProcessQueue(logger, agentGroupID)
	defer pq.Close()
	logger.Debug("Connected to process queue")

	velocity := velocity.NewAgent(os.Getenv("VELOCITY_URL"))

	// If this is in dev mode, we wait a little because
	// both services usually get restarted at the same time.
	if os.Getenv("DEV_SERVICES") == "true" {
		time.Sleep(2 * time.Second)
	}

	logger.Debug("Starting agent")
	agent := agent.New(pq, velocity, logger)
	err := agent.Start(ctx)
	if err != nil {
		panic(err)
	}
}
