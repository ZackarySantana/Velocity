package mock

import (
	"context"

	"github.com/zackarysantana/velocity/internal/service"
	"github.com/zackarysantana/velocity/src/entities/image"
	"github.com/zackarysantana/velocity/src/entities/job"
	"github.com/zackarysantana/velocity/src/entities/routine"
	"github.com/zackarysantana/velocity/src/entities/test"
)

func NewMockRepositoryManager[T comparable](idCreator service.IDCreator[T]) *service.RepositoryManager[T] {
	routines := make(map[T]*routine.Routine[T])
	jobs := make(map[T]*job.Job[T])
	images := make(map[T]*image.Image[T])
	tests := make(map[T]*test.Test[T])
	return &service.RepositoryManager[T]{
		Routine: &service.RoutineRepository[T]{
			Load: createLoad(routines),
			Put:  createPutForType(routines, idCreator),
		},
		Job: &service.JobRepository[T]{
			Load: createLoad(jobs),
			Put:  createPutForType(jobs, idCreator),
		},
		Image: &service.ImageRepository[T]{
			Load: createLoad(images),
			Put:  createPutForType(images, idCreator),
		},
		Test: &service.TestRepository[T]{
			Load: createLoad(tests),
			Put:  createPutForType(tests, idCreator),
		},
		WithTransaction: func(ctx context.Context, fn func(context.Context) error) error {
			beforeRoutines := make(map[T]*routine.Routine[T], len(routines))
			for k, v := range routines {
				beforeRoutines[k] = v
			}
			beforeJobs := make(map[T]*job.Job[T], len(jobs))
			for k, v := range jobs {
				beforeJobs[k] = v
			}
			beforeImages := make(map[T]*image.Image[T], len(images))
			for k, v := range images {
				beforeImages[k] = v
			}
			beforeTests := make(map[T]*test.Test[T], len(tests))
			for k, v := range tests {
				beforeTests[k] = v
			}
			err := fn(ctx)
			if err != nil {
				routines = beforeRoutines
				jobs = beforeJobs
				images = beforeImages
				tests = beforeTests
			}
			return err
		},
	}
}

func createLoad[T any, V comparable](items map[V]*T) func(context.Context, []V) ([]*T, error) {
	return func(ctx context.Context, v []V) ([]*T, error) {
		results := make([]*T, 0, len(v))
		for _, item := range v {
			if result, ok := items[item]; ok {
				results = append(results, result)
			}
		}
		return results, nil
	}
}

func createPutForType[T any, V comparable](items map[V]*T, idCreator service.IDCreator[V]) func(context.Context, []*T) ([]V, error) {
	return func(ctx context.Context, t []*T) ([]V, error) {
		ids := make([]V, 0, len(t))
		for _, item := range t {
			id := idCreator.Create()
			ids = append(ids, id)
			items[id] = item
		}
		return ids, nil
	}
}
