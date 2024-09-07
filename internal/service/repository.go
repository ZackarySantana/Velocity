package service

import (
	"context"

	"github.com/zackarysantana/velocity/src/entities/image"
	"github.com/zackarysantana/velocity/src/entities/job"
	"github.com/zackarysantana/velocity/src/entities/routine"
	"github.com/zackarysantana/velocity/src/entities/test"
)

type TypeRepository[ID any, DataType any] interface {
	// Load gets the data for the given IDs.
	Load(context.Context, []ID) ([]*DataType, error)
	// Put saves the data and returns the inserted IDs.
	Put(context.Context, []*DataType) ([]ID, error)
}

type RoutineRepository[ID any] TypeRepository[ID, routine.Routine[ID]]
type JobRepository[ID any] TypeRepository[ID, job.Job[ID]]
type ImageRepository[ID any] TypeRepository[ID, image.Image[ID]]
type TestRepository[ID any] TypeRepository[ID, test.Test[ID]]

type RepositoryManager[ID any] interface {
	Routine() RoutineRepository[ID]
	Job() JobRepository[ID]
	Image() ImageRepository[ID]
	Test() TestRepository[ID]
	WithTransaction(context.Context, func(context.Context) error) error
}

// NewRepositoryManager holds the given repositories and transaction function
// and returns/calls them when needed.
func NewRepositoryManager[ID any](
	routine RoutineRepository[ID],
	job JobRepository[ID],
	image ImageRepository[ID],
	test TestRepository[ID],
	withTransaction func(context.Context, func(context.Context) error) error,
) RepositoryManager[ID] {
	return &newRepositoryManagerShim[ID]{
		routine:         routine,
		job:             job,
		image:           image,
		test:            test,
		withTransaction: withTransaction,
	}
}

type newRepositoryManagerShim[ID any] struct {
	routine         RoutineRepository[ID]
	job             JobRepository[ID]
	image           ImageRepository[ID]
	test            TestRepository[ID]
	withTransaction func(context.Context, func(context.Context) error) error
}

func (r *newRepositoryManagerShim[ID]) Routine() RoutineRepository[ID] {
	return r.routine
}

func (r *newRepositoryManagerShim[ID]) Job() JobRepository[ID] {
	return r.job
}

func (r *newRepositoryManagerShim[ID]) Image() ImageRepository[ID] {
	return r.image
}

func (r *newRepositoryManagerShim[ID]) Test() TestRepository[ID] {
	return r.test
}

func (r *newRepositoryManagerShim[ID]) WithTransaction(ctx context.Context, fn func(context.Context) error) error {
	return r.withTransaction(ctx, fn)
}
