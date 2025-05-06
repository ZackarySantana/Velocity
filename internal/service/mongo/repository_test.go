package mongo

import (
	"context"
	"fmt"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zackarysantana/velocity/internal/mongotest"
	"github.com/zackarysantana/velocity/internal/service"
	"github.com/zackarysantana/velocity/internal/service/servicetest"
)

func TestMongoRepository(t *testing.T) {
	ctx := context.Background()

	client, cleanup, err := mongotest.CreateContainer(ctx)
	require.NoError(t, err)
	t.Cleanup(func() { cleanup(ctx) })

	var i *int32
	var starting int32 = 0
	i = &starting

	repoGen := func() service.RepositoryManager[any] {
		return NewRepositoryManager[any](client, fmt.Sprintf("test_%d", atomic.AddInt32(i, 1)))
	}

	servicetest.TestRepository(t, repoGen, NewIDCreator[any]())
}
