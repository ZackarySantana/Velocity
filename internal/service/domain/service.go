package domain

import (
	"context"
	"log/slog"

	"github.com/zackarysantana/velocity/internal/service"
	"github.com/zackarysantana/velocity/src/entities"
)

type Service[ID any] struct {
	idCreator  service.IDCreator[ID]
	repository service.RepositoryManager[ID]

	testQueue service.PriorityQueue[ID, ID]

	logger *slog.Logger
}

func NewService[T any](idCreator service.IDCreator[T], repository service.RepositoryManager[T], pq service.PriorityQueue[T, T], logger *slog.Logger) service.Service[T] {
	return &Service[T]{repository: repository, testQueue: pq, idCreator: idCreator, logger: logger}
}

func (s *Service[T]) StartRoutine(ctx context.Context, ec *entities.ConfigEntity[T], name string) error {
	return s.repository.WithTransaction(ctx, func(ctx context.Context) error {
		ids, err := s.repository.Test().Put(ctx, ec.Tests)
		if err != nil {
			return err
		}
		for i := range ec.Tests {
			ec.Tests[i].Id = ids[i]
		}

		ids, err = s.repository.Image().Put(ctx, ec.Images)
		if err != nil {
			return err
		}
		for i := range ec.Images {
			ec.Images[i].Id = ids[i]
		}

		ids, err = s.repository.Job().Put(ctx, ec.Jobs)
		if err != nil {
			return err
		}
		for i := range ec.Jobs {
			ec.Jobs[i].Id = ids[i]
		}

		ids, err = s.repository.Routine().Put(ctx, ec.Routines)
		if err != nil {
			return err
		}
		for i := range ec.Routines {
			ec.Routines[i].Id = ids[i]
		}

		tests := make([]service.PriorityQueueItem[T], len(ec.Tests))
		for i, t := range ec.Tests {
			tests[i] = service.PriorityQueueItem[T]{Priority: 1, Payload: t.Id}
		}

		return s.testQueue.Push(ctx, "tests", tests...)
	})
}
