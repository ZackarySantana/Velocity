package domain

import (
	"context"
	"log/slog"

	"github.com/samber/oops"
	"github.com/zackarysantana/velocity/internal/service"
	"github.com/zackarysantana/velocity/src/entities"
	"github.com/zackarysantana/velocity/src/entities/routine"
)

type Service[ID any] struct {
	idCreator  service.IDCreator[ID]
	repository service.RepositoryManager[ID]

	testQueue service.PriorityQueue[ID, ID]

	logger *slog.Logger
}

func NewService[T any](repository service.RepositoryManager[T], pq service.PriorityQueue[T, T], idCreator service.IDCreator[T], logger *slog.Logger) service.Service[T] {
	return &Service[T]{repository: repository, testQueue: pq, idCreator: idCreator, logger: logger}
}

func (s *Service[T]) StartRoutine(ctx context.Context, ec *entities.ConfigEntity[T], name string) error {
	return s.repository.WithTransaction(ctx, func(ctx context.Context) error {
		_, err := s.repository.Test().Put(ctx, ec.Tests)
		if err != nil {
			return err
		}

		_, err = s.repository.Image().Put(ctx, ec.Images)
		if err != nil {
			return err
		}

		_, err = s.repository.Job().Put(ctx, ec.Jobs)
		if err != nil {
			return err
		}

		for _, r := range ec.Routines {
			_, err := s.repository.Routine().Put(ctx, []*routine.Routine[T]{r})
			if err != nil {
				return err
			}
		}

		// testIds := make([][]byte, len(ec.Tests))
		tests := make([]service.PriorityQueueItem[T], len(ec.Tests))
		for i, t := range ec.Tests {
			id, err := s.idCreator.Read(t.Id)
			if err != nil {
				return oops.Wrapf(err, "id is invalid")
			}
			tests[i] = service.PriorityQueueItem[T]{Priority: 1, Payload: id}
		}

		return s.testQueue.Push(ctx, "tests", tests...)
	})
}
