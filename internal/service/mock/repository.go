package mock

import (
	"context"

	"github.com/zackarysantana/velocity/internal/service"
	"github.com/zackarysantana/velocity/src/entities/image"
	"github.com/zackarysantana/velocity/src/entities/job"
	"github.com/zackarysantana/velocity/src/entities/routine"
	"github.com/zackarysantana/velocity/src/entities/test"
)

func NewRepositoryManager[ID comparable](idCreator service.IDCreator[ID]) service.RepositoryManager[ID] {
	routines := make(map[ID]*routine.Routine[ID])
	jobs := make(map[ID]*job.Job[ID])
	images := make(map[ID]*image.Image[ID])
	tests := make(map[ID]*test.Test[ID])
	return service.NewRepositoryManager[ID](
		newTypeRepository[ID, routine.Routine[ID]](idCreator, routines),
		newTypeRepository[ID, job.Job[ID]](idCreator, jobs),
		newTypeRepository[ID, image.Image[ID]](idCreator, images),
		newTypeRepository[ID, test.Test[ID]](idCreator, tests),
		func(ctx context.Context, fn func(context.Context) error) error {
			beforeRoutines := make(map[ID]*routine.Routine[ID], len(routines))
			for k, v := range routines {
				beforeRoutines[k] = v
			}
			beforeJobs := make(map[ID]*job.Job[ID], len(jobs))
			for k, v := range jobs {
				beforeJobs[k] = v
			}
			beforeImages := make(map[ID]*image.Image[ID], len(images))
			for k, v := range images {
				beforeImages[k] = v
			}
			beforeTests := make(map[ID]*test.Test[ID], len(tests))
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
	)
}

func newTypeRepository[ID comparable, DataType any](idCreator service.IDCreator[ID], data map[ID]*DataType) service.TypeRepository[ID, DataType] {
	return &typeRepository[ID, DataType]{idCreator: idCreator, data: data}
}

type typeRepository[ID comparable, DataType any] struct {
	idCreator service.IDCreator[ID]
	data      map[ID]*DataType
}

func (r *typeRepository[ID, DataType]) Load(ctx context.Context, ids []ID) ([]*DataType, error) {
	results := make([]*DataType, 0, len(ids))
	for _, id := range ids {
		if result, ok := r.data[id]; ok {
			results = append(results, result)
		}
	}
	return results, nil
}

func (r *typeRepository[ID, DataType]) Put(ctx context.Context, data []*DataType) ([]ID, error) {
	ids := make([]ID, 0, len(data))
	for _, item := range data {
		id := r.idCreator.Create()
		ids = append(ids, id)
		r.data[id] = item
	}
	return ids, nil
}
