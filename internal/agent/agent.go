package agent

import (
	"context"
	"fmt"

	"github.com/zackarysantana/velocity/internal/service"
	"github.com/zackarysantana/velocity/src/velocity"
)

type agent struct {
	processQueue service.ProcessQueue

	client *velocity.AgentClient
}

func New(processQueue service.ProcessQueue, client *velocity.AgentClient) *agent {
	return &agent{processQueue: processQueue, client: client}
}

func (a *agent) Start(ctx context.Context) error {
	fmt.Println("Pinging server...")
	_, err := a.client.Health()
	if err != nil {
		return err
	}
	fmt.Println("Pinged server")

	err = a.processQueue.Consume(ctx, "tests", func(data []byte) (bool, error) {
		id := string(data)
		fmt.Println("Received test:", id)
		fmt.Println("Attempting to get test...")
		_, res, err := a.client.GetTest(id)
		if err != nil {
			return false, err
		}
		fmt.Println("Got test", res)
		return true, nil
	})
	if err != nil {
		if err == context.Canceled {
			return nil
		}
	}
	return err
}
