package domain

import (
	"context"
	"log/slog"

	"github.com/zackarysantana/velocity/internal/service"
	"github.com/zackarysantana/velocity/src/entities"
	"github.com/zackarysantana/velocity/src/entities/routine"
)

type Service struct {
	repository *service.Repository
	pq         service.ProcessQueue
	logger     *slog.Logger
}

func NewService(repository *service.Repository, pq service.ProcessQueue, logger *slog.Logger) service.Service {
	return &Service{repository: repository, pq: pq, logger: logger}
}

func (s *Service) StartRoutine(ctx context.Context, ec *entities.ConfigEntity, name string) error {
	return s.repository.WithTransaction(ctx, func(ctx context.Context) error {
		err := s.repository.Test.Put(ctx, ec.Tests)
		if err != nil {
			return err
		}

		err = s.repository.Image.Put(ctx, ec.Images)
		if err != nil {
			return err
		}

		err = s.repository.Job.Put(ctx, ec.Jobs)
		if err != nil {
			return err
		}

		for _, r := range ec.Routines {
			err := s.repository.Routine.Put(ctx, []*routine.Routine{r})
			if err != nil {
				return err
			}
		}

		testIds := make([][]byte, len(ec.Tests))
		for i, t := range ec.Tests {
			testIds[i] = []byte(t.Id)
		}

		return s.pq.Write(ctx, "tests", testIds...)
	})
}
