package domain

import (
	"context"

	"github.com/zackarysantana/velocity/internal/service"
	"github.com/zackarysantana/velocity/src/entities"
	"github.com/zackarysantana/velocity/src/entities/routine"
)

type Service struct {
	db *service.Repository
}

func NewService(db *service.Repository) service.Service {
	return &Service{db: db}
}

func (s *Service) StartRoutine(ctx context.Context, ec *entities.ConfigEntity, name string) error {
	return s.db.WithTransaction(ctx, func(ctx context.Context) error {
		err := s.db.Test.Put(ctx, ec.Tests)
		if err != nil {
			return err
		}

		err = s.db.Image.Put(ctx, ec.Images)
		if err != nil {
			return err
		}

		err = s.db.Job.Put(ctx, ec.Jobs)
		if err != nil {
			return err
		}

		for _, r := range ec.Routines {
			if r.Name != name {
				continue
			}
			err := s.db.Routine.Put(ctx, []*routine.Routine{r})
			if err != nil {
				return err
			}
			break
		}

		return nil
	})
}
