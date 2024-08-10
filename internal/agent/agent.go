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
	resp, err := a.client.Health()
	if err != nil {
		return err
	}
	fmt.Println(*resp)

	err = a.processQueue.Consume(ctx, "tests", func(data []byte) error {
		// id := string(data)

		return nil
	})
	if err != nil {
		if err == context.Canceled {
			return nil
		}
	}
	return err
}
