package servicetest

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zackarysantana/velocity/internal/service"
)

func TestPriorityQueue(t *testing.T, pqGen func() service.PriorityQueue[any, string]) {
	ctx := context.Background()

	for tName, tFunc := range map[string]func(*testing.T, service.PriorityQueue[any, string]){
		"SingleItem": func(t *testing.T, pq service.PriorityQueue[any, string]) {
			err := pq.Push(ctx, "test", service.PriorityQueueItem[string]{Priority: 1, Payload: "test"})
			require.NoError(t, err)

			item, err := pq.Pop(ctx, "test")
			require.NoError(t, err)
			assert.Equal(t, "test", item.Payload)

			item, err = pq.Pop(ctx, "test")
			assert.Equal(t, service.ErrEmptyQueue, err)
		},
		"MultiItem": func(t *testing.T, pq service.PriorityQueue[any, string]) {
			err := pq.Push(ctx, "test", service.PriorityQueueItem[string]{Priority: 1, Payload: "test1"}, service.PriorityQueueItem[string]{Priority: 2, Payload: "test2"})
			require.NoError(t, err)

			item, err := pq.Pop(ctx, "test")
			require.NoError(t, err)
			require.Equal(t, "test2", item.Payload)

			item, err = pq.Pop(ctx, "test")
			require.NoError(t, err)
			require.Equal(t, "test1", item.Payload)

			item, err = pq.Pop(ctx, "test")
			assert.Equal(t, service.ErrEmptyQueue, err)
		},
		"MutliQueue": func(t *testing.T, pq service.PriorityQueue[any, string]) {
			err := pq.Push(ctx, "test1", service.PriorityQueueItem[string]{Priority: 1, Payload: "test1"})
			require.NoError(t, err)

			err = pq.Push(ctx, "test2", service.PriorityQueueItem[string]{Priority: 2, Payload: "test2"})
			require.NoError(t, err)

			item, err := pq.Pop(ctx, "test2")
			require.NoError(t, err)
			require.Equal(t, "test2", item.Payload)

			item, err = pq.Pop(ctx, "test2")
			assert.Equal(t, service.ErrEmptyQueue, err)

			item, err = pq.Pop(ctx, "test1")
			require.NoError(t, err)
			require.Equal(t, "test1", item.Payload)

			item, err = pq.Pop(ctx, "test1")
			assert.Equal(t, service.ErrEmptyQueue, err)
		},
	} {
		t.Run(tName, func(t *testing.T) {
			tFunc(t, pqGen())
		})
	}
}
