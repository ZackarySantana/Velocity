package servicetest

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zackarysantana/velocity/internal/service"
	"github.com/zackarysantana/velocity/src/entities/image"
	"github.com/zackarysantana/velocity/src/entities/job"
	"github.com/zackarysantana/velocity/src/entities/routine"
	"github.com/zackarysantana/velocity/src/entities/test"
)

func TestRepository(t *testing.T, repoGen func() service.RepositoryManager[any], idCreator service.IDCreator[any]) {
	ctx := context.Background()

	for name, testFunc := range map[string]func(t *testing.T, rm service.RepositoryManager[any]){
		"RoutinePutAndLoad": func(t *testing.T, rm service.RepositoryManager[any]) {
			routineName := "routine1"
			r := &routine.Routine[any]{Id: idCreator.Create(), Name: routineName}
			ids, err := rm.Routine().Put(ctx, []*routine.Routine[any]{r})
			require.NoError(t, err)
			require.Len(t, ids, 1)

			results, err := rm.Routine().Load(ctx, ids)
			require.NoError(t, err)
			require.Len(t, results, 1)
			assert.Equal(t, routineName, results[0].Name)
		},
		"JobPutAndLoad": func(t *testing.T, rm service.RepositoryManager[any]) {
			jobName := "job1"
			j := &job.Job[any]{Id: idCreator.Create(), Name: jobName}
			ids, err := rm.Job().Put(ctx, []*job.Job[any]{j})
			require.NoError(t, err)
			require.Len(t, ids, 1)

			results, err := rm.Job().Load(ctx, ids)
			require.NoError(t, err)
			require.Len(t, results, 1)
			assert.Equal(t, jobName, results[0].Name)
		},
		"ImagePutAndLoad": func(t *testing.T, rm service.RepositoryManager[any]) {
			imageName := "image1"
			i := &image.Image[any]{Id: idCreator.Create(), Name: imageName}
			ids, err := rm.Image().Put(ctx, []*image.Image[any]{i})
			require.NoError(t, err)
			require.Len(t, ids, 1)

			results, err := rm.Image().Load(ctx, ids)
			require.NoError(t, err)
			require.Len(t, results, 1)
			assert.Equal(t, imageName, results[0].Name)
		},
		"TestPutAndLoad": func(t *testing.T, rm service.RepositoryManager[any]) {
			testName := "test1"
			ts := &test.Test[any]{Id: idCreator.Create(), Name: testName}
			ids, err := rm.Test().Put(ctx, []*test.Test[any]{ts})
			require.NoError(t, err)
			require.Len(t, ids, 1)

			results, err := rm.Test().Load(ctx, ids)
			require.NoError(t, err)
			require.Len(t, results, 1)
			assert.Equal(t, testName, results[0].Name)
		},
		"TransactionRollbackOnError": func(t *testing.T, rm service.RepositoryManager[any]) {
			routineName := "txRoutine"
			r := &routine.Routine[any]{Id: idCreator.Create(), Name: routineName}
			var insertedRoutineID any

			jobName := "txJob"
			j := &job.Job[any]{Id: idCreator.Create(), Name: jobName}
			var insertedJobID any

			imageName := "txImage"
			i := &image.Image[any]{Id: idCreator.Create(), Name: imageName}
			var insertedImageID any

			testName := "txTest"
			ts := &test.Test[any]{Id: idCreator.Create(), Name: testName}
			var insertedTestID any

			err := rm.WithTransaction(ctx, func(txCtx context.Context) error {
				ids, err := rm.Routine().Put(txCtx, []*routine.Routine[any]{r})
				require.NoError(t, err)
				require.Len(t, ids, 1)
				insertedRoutineID = ids[0]

				ids, err = rm.Job().Put(txCtx, []*job.Job[any]{j})
				require.NoError(t, err)
				require.Len(t, ids, 1)
				insertedJobID = ids[0]

				ids, err = rm.Image().Put(txCtx, []*image.Image[any]{i})
				require.NoError(t, err)
				require.Len(t, ids, 1)
				insertedImageID = ids[0]

				ids, err = rm.Test().Put(txCtx, []*test.Test[any]{ts})
				require.NoError(t, err)
				require.Len(t, ids, 1)
				insertedTestID = ids[0]

				return assert.AnError
			})
			require.Error(t, err)

			routineResults, err := rm.Routine().Load(ctx, []any{insertedRoutineID})
			require.NoError(t, err)
			assert.Len(t, routineResults, 0)

			jobResults, err := rm.Job().Load(ctx, []any{insertedJobID})
			require.NoError(t, err)
			assert.Len(t, jobResults, 0)

			imageResults, err := rm.Image().Load(ctx, []any{insertedImageID})
			require.NoError(t, err)
			assert.Len(t, imageResults, 0)

			testResults, err := rm.Test().Load(ctx, []any{insertedTestID})
			require.NoError(t, err)
			assert.Len(t, testResults, 0)
		},
		"TransactionCommitOnSuccess": func(t *testing.T, rm service.RepositoryManager[any]) {
			routineName := "committedRoutine"
			r := &routine.Routine[any]{Id: idCreator.Create(), Name: routineName}
			var insertedRoutineID any

			jobName := "committedJob"
			j := &job.Job[any]{Id: idCreator.Create(), Name: jobName}
			var insertedJobID any

			imageName := "committedImage"
			i := &image.Image[any]{Id: idCreator.Create(), Name: imageName}
			var insertedImageID any

			testName := "committedTest"
			ts := &test.Test[any]{Id: idCreator.Create(), Name: testName}
			var insertedTestID any

			err := rm.WithTransaction(ctx, func(txCtx context.Context) error {
				ids, err := rm.Routine().Put(txCtx, []*routine.Routine[any]{r})
				require.NoError(t, err)
				require.Len(t, ids, 1)
				insertedRoutineID = ids[0]

				ids, err = rm.Job().Put(txCtx, []*job.Job[any]{j})
				require.NoError(t, err)
				require.Len(t, ids, 1)
				insertedJobID = ids[0]

				ids, err = rm.Image().Put(txCtx, []*image.Image[any]{i})
				require.NoError(t, err)
				require.Len(t, ids, 1)
				insertedImageID = ids[0]

				ids, err = rm.Test().Put(txCtx, []*test.Test[any]{ts})
				require.NoError(t, err)
				require.Len(t, ids, 1)
				insertedTestID = ids[0]

				return nil
			})
			require.NoError(t, err)

			routineResults, err := rm.Routine().Load(ctx, []any{insertedRoutineID})
			require.NoError(t, err)
			require.Len(t, routineResults, 1)
			assert.Equal(t, routineName, routineResults[0].Name)

			jobResults, err := rm.Job().Load(ctx, []any{insertedJobID})
			require.NoError(t, err)
			require.Len(t, jobResults, 1)
			assert.Equal(t, jobName, jobResults[0].Name)

			imageResults, err := rm.Image().Load(ctx, []any{insertedImageID})
			require.NoError(t, err)
			require.Len(t, imageResults, 1)
			assert.Equal(t, imageName, imageResults[0].Name)

			testResults, err := rm.Test().Load(ctx, []any{insertedTestID})
			require.NoError(t, err)
			require.Len(t, testResults, 1)
			assert.Equal(t, testName, testResults[0].Name)
		},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			testFunc(t, repoGen())
		})
	}
}
