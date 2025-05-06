package mongo

import (
	"context"
	"fmt"
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
	defer cleanup(ctx)

	i := 0
	repoGen := func() service.RepositoryManager[any] {
		i++
		return NewRepositoryManager[any](client, fmt.Sprintf("test_%d", i))
	}

	servicetest.TestRepository(t, repoGen)
}
