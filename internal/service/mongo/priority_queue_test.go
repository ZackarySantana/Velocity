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

func TestMongoPriorityQueue(t *testing.T) {
	ctx := context.Background()

	client, cleanup, err := mongotest.CreateContainer(ctx)
	require.NoError(t, err)
	t.Cleanup(func() { cleanup(ctx) })

	i := 0
	pqGen := func() service.PriorityQueue[any, string] {
		i++
		return NewPriorityQueue[any, string](client, NewIDCreator[any](), fmt.Sprintf("test_%d", i))
	}

	servicetest.TestPriorityQueue(t, pqGen)
}
