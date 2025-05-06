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
	return service.NewRepositoryManager(
		newTypeRepository(idCreator, routines),
		newTypeRepository(idCreator, jobs),
		newTypeRepository(idCreator, images),
		newTypeRepository(idCreator, tests),
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
				// We have to clear and reinsert because the maps aren't pointers.
				for k := range routines {
					delete(routines, k)
				}
				for k := range jobs {
					delete(jobs, k)
				}
				for k := range images {
					delete(images, k)
				}
				for k := range tests {
					delete(tests, k)
				}
				for k, v := range beforeRoutines {
					routines[k] = v
				}
				for k, v := range beforeJobs {
					jobs[k] = v
				}
				for k, v := range beforeImages {
					images[k] = v
				}
				for k, v := range beforeTests {
					tests[k] = v
				}
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
