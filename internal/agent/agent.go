package agent

import (
	"context"
	"log/slog"
	"time"

	"github.com/samber/oops"
	"github.com/zackarysantana/velocity/internal/service"
	"github.com/zackarysantana/velocity/src/velocity"
)

type agent struct {
	priorityQueue service.PriorityQueue[any, any]

	client *velocity.AgentClient

	logger *slog.Logger
}

func New(priorityQueue service.PriorityQueue[any, any], client *velocity.AgentClient, logger *slog.Logger) *agent {
	return &agent{priorityQueue: priorityQueue, client: client, logger: logger}
}

func (a *agent) Start(ctx context.Context) error {
	a.logger.Debug("Pinging server...")
	_, err := a.client.Health(ctx)
	if err != nil {
		return oops.Wrapf(err, "failed to ping server")
	}
	a.logger.Debug("Pinged server")

	for {
		select {
		case <-ctx.Done():
			a.logger.Debug("Context cancelled. Exiting")
			return ctx.Err()
		default:
			item, err := a.priorityQueue.Pop(ctx, "tests")
			if err != nil {
				if err == service.ErrEmptyQueue {
					a.logger.Debug("No tests found. Waiting for new tests")
					// TODO: Proper backoff strategy for all continues
					time.Sleep(time.Second)
					continue
				}
				return oops.Wrapf(err, "failed to pop from queue")
			}
			testID, ok := item.Payload.(string)
			if !ok {
				return oops.Errorf("could not convert %v to test ID", item.Payload)
			}
			a.logger.Debug("Received test", "id", testID)

			resp, data, err := a.client.GetTest(ctx, testID)
			if err != nil {
				if resp.StatusCode == 404 {
					a.logger.Debug("Test not found. Skipping", "id", testID)
					continue
				}
				return oops.Wrapf(err, "failed to get test")
			}
			a.logger.Debug("Found test", "test", data)

			err = a.priorityQueue.MarkAsDone(ctx, "tests", item.Id)
			if err != nil {
				return oops.Wrapf(err, "failed to mark test as done")
			}
		}
	}
}
