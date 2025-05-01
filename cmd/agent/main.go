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
	pq := internal.GetPriorityQueue[any, any](logger)
	defer pq.Close()
	logger.Debug("Connected to process queue")

	velocity := velocity.NewAgentClient(os.Getenv("VELOCITY_URL"))

	// If this is in dev mode, we wait a little because
	// both services usually get restarted at the same time.
	if os.Getenv("DEV_SERVICES") == "true" {
		time.Sleep(2 * time.Second)
	}

	logger.Debug("Starting agent")
	agent := agent.New(pq, velocity, logger)

	// TODO: make this reset when enough time has passed.
	errors := []error{}
	for {
		if len(errors) == 4 {
			logger.Error("Agent failed 5 times. Exiting")
			os.Exit(1)
		}
		if len(errors) > 0 {
			logger.Error("Agent failed. Restarting in 5 seconds", "errors", errors)
			time.Sleep(5 * time.Second)
		}

		err := agent.Start(ctx)
		if err != nil {
			logger.Error("Agent failed", "error", err)
			errors = append(errors, err)
		}
	}
}
