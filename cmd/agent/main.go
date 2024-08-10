package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/zackarysantana/velocity/internal/agent"
	"github.com/zackarysantana/velocity/internal/service/kafka"
	"github.com/zackarysantana/velocity/src/velocity"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file", err)
	}

	pq, err := kafka.NewKafkaQueue(os.Getenv("KAFKA_USERNAME"), os.Getenv("KAFKA_PASSWORD"), os.Getenv("KAFKA_BROKER"), "agent")
	defer pq.Close()
	if err != nil {
		panic(err)
	}

	velocity := velocity.NewAgent(os.Getenv("VELOCITY_URL"))

	// If this is in dev mode, we wait a little because
	// both services usually get restarted at the same time.
	if os.Getenv("DEV_MODE") == "true" {
		time.Sleep(2 * time.Second)
	}

	agent := agent.New(pq, velocity)
	err = agent.Start(context.Background())
	if err != nil {
		panic(err)
	}
}
