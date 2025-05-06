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

func TestRepository(t *testing.T, repoGen func() service.RepositoryManager[any]) {
	ctx := context.Background()

	for name, testFunc := range map[string]func(t *testing.T, rm service.RepositoryManager[any]){
		"RoutinePutAndLoad": func(t *testing.T, rm service.RepositoryManager[any]) {
			r := &routine.Routine[any]{Name: "routine1"}
			ids, err := rm.Routine().Put(ctx, []*routine.Routine[any]{r})
			require.NoError(t, err)
			require.Len(t, ids, 1)

			results, err := rm.Routine().Load(ctx, ids)
			require.NoError(t, err)
			require.Len(t, results, 1)
			assert.Equal(t, "routine1", results[0].Name)
		},
		"JobPutAndLoad": func(t *testing.T, rm service.RepositoryManager[any]) {
			j := &job.Job[any]{Name: "job1"}
			ids, err := rm.Job().Put(ctx, []*job.Job[any]{j})
			require.NoError(t, err)
			require.Len(t, ids, 1)

			results, err := rm.Job().Load(ctx, ids)
			require.NoError(t, err)
			require.Len(t, results, 1)
			assert.Equal(t, "job1", results[0].Name)
		},
		"ImagePutAndLoad": func(t *testing.T, rm service.RepositoryManager[any]) {
			img := &image.Image[any]{Name: "image1"}
			ids, err := rm.Image().Put(ctx, []*image.Image[any]{img})
			require.NoError(t, err)
			require.Len(t, ids, 1)

			results, err := rm.Image().Load(ctx, ids)
			require.NoError(t, err)
			require.Len(t, results, 1)
			assert.Equal(t, "image1", results[0].Name)
		},
		"TestPutAndLoad": func(t *testing.T, rm service.RepositoryManager[any]) {
			ts := &test.Test[any]{Name: "test1"}
			ids, err := rm.Test().Put(ctx, []*test.Test[any]{ts})
			require.NoError(t, err)
			require.Len(t, ids, 1)

			results, err := rm.Test().Load(ctx, ids)
			require.NoError(t, err)
			require.Len(t, results, 1)
			assert.Equal(t, "test1", results[0].Name)
		},
		"TransactionRollbackOnError": func(t *testing.T, rm service.RepositoryManager[any]) {
			r := &routine.Routine[any]{Name: "txRoutine"}
			var insertedID any

			err := rm.WithTransaction(ctx, func(txCtx context.Context) error {
				ids, err := rm.Routine().Put(txCtx, []*routine.Routine[any]{r})
				require.NoError(t, err)
				insertedID = ids[0]
				return assert.AnError
			})
			require.Error(t, err)

			results, err := rm.Routine().Load(ctx, []any{insertedID})
			require.NoError(t, err)
			assert.Len(t, results, 0)
		},
		"TransactionCommitOnSuccess": func(t *testing.T, rm service.RepositoryManager[any]) {
			r := &routine.Routine[any]{Name: "committedRoutine"}
			var insertedID any

			err := rm.WithTransaction(ctx, func(txCtx context.Context) error {
				ids, err := rm.Routine().Put(txCtx, []*routine.Routine[any]{r})
				require.NoError(t, err)
				insertedID = ids[0]
				return nil
			})
			require.NoError(t, err)

			results, err := rm.Routine().Load(ctx, []any{insertedID})
			require.NoError(t, err)
			require.Len(t, results, 1)
			assert.Equal(t, "committedRoutine", results[0].Name)
		},
	} {
		t.Run(name, func(t *testing.T) {
			testFunc(t, repoGen())
		})
	}
}
