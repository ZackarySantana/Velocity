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
	_, err := a.client.Health()
	if err != nil {
		return err
	}

	err = a.processQueue.Consume(ctx, "tests", func(data []byte) error {
		id := string(data)
		_, res, err := a.client.GetTest(id)
		fmt.Println("Trying to find:", id)
		if err != nil {
			return err
		}
		fmt.Println("for: ", res)
		return nil
	})
	if err != nil {
		if err == context.Canceled {
			return nil
		}
	}
	return err
}
