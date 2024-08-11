package agent

import (
	"context"
	"log/slog"

	"github.com/zackarysantana/velocity/internal/service"
	"github.com/zackarysantana/velocity/src/velocity"
)

type agent struct {
	processQueue service.ProcessQueue

	client *velocity.AgentClient

	logger *slog.Logger
}

func New(processQueue service.ProcessQueue, client *velocity.AgentClient, logger *slog.Logger) *agent {
	return &agent{processQueue: processQueue, client: client, logger: logger}
}

func (a *agent) Start(ctx context.Context) error {
	a.logger.Debug("Pinging server...")
	_, err := a.client.Health()
	if err != nil {
		return err
	}
	a.logger.Debug("Pinged server")

	err = a.processQueue.Consume(ctx, "tests", func(data []byte) (bool, error) {
		id := string(data)
		a.logger.Debug("Received test", "id", id)
		a.logger.Debug("Attempting to get test...")
		_, res, err := a.client.GetTest(id)
		if err != nil {
			return false, err
		}
		a.logger.Debug("Got test", "test", res)
		return true, nil
	})
	if err != nil {
		if err == context.Canceled {
			return nil
		}
	}
	return err
}
