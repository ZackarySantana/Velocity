package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/zackarysantana/velocity/internal/agent"
	"github.com/zackarysantana/velocity/internal/service/kafka"
	"github.com/zackarysantana/velocity/src/velocity"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file", err)
	}

	pq, err := kafka.NewKafkaQueue(os.Getenv("KAFKA_USERNAME"), os.Getenv("KAFKA_PASSWORD"), os.Getenv("KAFKA_BROKER"), "test-group-12")
	defer pq.Close()
	if err != nil {
		panic(err)
	}

	velocity := velocity.NewAgent(os.Getenv("VELOCITY_URL"))

	agent := agent.New(pq, velocity)
	err = agent.Start(context.Background())
	if err != nil {
		panic(err)
	}
}
