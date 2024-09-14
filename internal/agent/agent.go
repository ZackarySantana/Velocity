package agent

import (
	"context"
	"log/slog"

	"github.com/samber/oops"
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
	_, err := a.client.Health(ctx)
	if err != nil {
		return oops.Wrapf(err, "failed to ping server")
	}
	a.logger.Debug("Pinged server")

	err = a.processQueue.Consume(ctx, "tests", func(idMsg []byte) (bool, error) {
		id := string(idMsg)
		a.logger.Debug("Received test", "id", id)
		resp, data, err := a.client.GetTest(ctx, id)
		if err != nil {
			if resp.StatusCode == 404 {
				a.logger.Debug("Test not found. Skipping", "id", id)
				return true, nil
			}
			return false, err
		}
		a.logger.Debug("Found test", "test", data)
		return true, nil
	})
	if err == context.Canceled {
		return nil
	}
	return err
}
