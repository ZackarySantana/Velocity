package domain

import (
	"context"
	"log/slog"

	"github.com/samber/oops"
	"github.com/zackarysantana/velocity/internal/service"
	"github.com/zackarysantana/velocity/src/entities"
	"github.com/zackarysantana/velocity/src/entities/routine"
)

type Service[T any] struct {
	repository *service.RepositoryManager[T]
	pq         service.ProcessQueue
	idCreator  service.IDCreator[T]
	logger     *slog.Logger
}

func NewService[T any](repository *service.RepositoryManager[T], pq service.ProcessQueue, idCreator service.IDCreator[T], logger *slog.Logger) service.Service[T] {
	return &Service[T]{repository: repository, pq: pq, idCreator: idCreator, logger: logger}
}

func (s *Service[T]) StartRoutine(ctx context.Context, ec *entities.ConfigEntity[T], name string) error {
	return s.repository.WithTransaction(ctx, func(ctx context.Context) error {
		_, err := s.repository.Test.Put(ctx, ec.Tests)
		if err != nil {
			return err
		}

		_, err = s.repository.Image.Put(ctx, ec.Images)
		if err != nil {
			return err
		}

		_, err = s.repository.Job.Put(ctx, ec.Jobs)
		if err != nil {
			return err
		}

		for _, r := range ec.Routines {
			_, err := s.repository.Routine.Put(ctx, []*routine.Routine[T]{r})
			if err != nil {
				return err
			}
		}

		testIds := make([][]byte, len(ec.Tests))
		for i, t := range ec.Tests {
			id, err := s.idCreator.Read(t.Id)
			if err != nil {
				return oops.Wrapf(err, "id is invalid")
			}
			testIds[i] = []byte(s.idCreator.String(id))
		}

		return s.pq.Write(ctx, "tests", testIds...)
	})
}
