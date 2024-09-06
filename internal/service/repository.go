package service

import (
	"context"

	"github.com/zackarysantana/velocity/src/entities/image"
	"github.com/zackarysantana/velocity/src/entities/job"
	"github.com/zackarysantana/velocity/src/entities/routine"
	"github.com/zackarysantana/velocity/src/entities/test"
)

type RepositoryManager[T any] struct {
	Routine *RoutineRepository[T]
	Job     *JobRepository[T]
	Image   *ImageRepository[T]
	Test    *TestRepository[T]

	WithTransaction func(context.Context, func(context.Context) error) error
}

type RoutineRepository[T any] struct {
	Load func(context.Context, []T) ([]*routine.Routine[T], error)
	Put  func(context.Context, []*routine.Routine[T]) ([]T, error)
}

type JobRepository[T any] struct {
	Load func(context.Context, []T) ([]*job.Job[T], error)
	Put  func(context.Context, []*job.Job[T]) ([]T, error)
}

type ImageRepository[T any] struct {
	Load func(context.Context, []T) ([]*image.Image[T], error)
	Put  func(context.Context, []*image.Image[T]) ([]T, error)
}

type TestRepository[T any] struct {
	Load func(context.Context, []T) ([]*test.Test[T], error)
	Put  func(context.Context, []*test.Test[T]) ([]T, error)
}
