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
			assert.ErrorIs(t, service.ErrEmptyQueue, err)
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
			assert.ErrorIs(t, service.ErrEmptyQueue, err)
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
			assert.ErrorIs(t, service.ErrEmptyQueue, err)

			item, err = pq.Pop(ctx, "test1")
			require.NoError(t, err)
			require.Equal(t, "test1", item.Payload)

			item, err = pq.Pop(ctx, "test1")
			assert.ErrorIs(t, service.ErrEmptyQueue, err)
		},
		"UnfinishedItems": func(t *testing.T, pq service.PriorityQueue[any, string]) {
			err := pq.Push(ctx, "test", service.PriorityQueueItem[string]{Priority: 1, Payload: "test1"}, service.PriorityQueueItem[string]{Priority: 2, Payload: "test2"})
			require.NoError(t, err)

			unfinishedItems, err := pq.UnfinishedItems(ctx, "test")
			require.NoError(t, err)
			require.Len(t, unfinishedItems, 0)

			// Popping should add one to the unfinished items.
			item, err := pq.Pop(ctx, "test")
			require.NoError(t, err)
			require.Equal(t, "test2", item.Payload)

			unfinishedItems, err = pq.UnfinishedItems(ctx, "test")
			require.NoError(t, err)
			require.Len(t, unfinishedItems, 1)
			foundTest1 := false
			foundTest2 := false
			for _, item := range unfinishedItems {
				if item.Payload == "test1" {
					foundTest1 = true
				}
				if item.Payload == "test2" {
					foundTest2 = true
				}
			}
			assert.False(t, foundTest1)
			assert.True(t, foundTest2)

			// Popping again should add one to the unfinished items.
			item, err = pq.Pop(ctx, "test")
			require.NoError(t, err)
			require.Equal(t, "test1", item.Payload)

			unfinishedItems, err = pq.UnfinishedItems(ctx, "test")
			require.NoError(t, err)
			require.Len(t, unfinishedItems, 2)
			foundTest1 = false
			foundTest2 = false
			for _, item := range unfinishedItems {
				if item.Payload == "test1" {
					foundTest1 = true
				}
				if item.Payload == "test2" {
					foundTest2 = true
				}
			}
			assert.True(t, foundTest1)
			assert.True(t, foundTest2)

			// Marking the item as done should remove it from the unfinished items.
			err = pq.MarkAsDone(ctx, "test", item.Id)
			require.NoError(t, err)

			unfinishedItems, err = pq.UnfinishedItems(ctx, "test")
			require.NoError(t, err)
			require.Len(t, unfinishedItems, 1)
			foundTest1 = false
			foundTest2 = false
			for _, item := range unfinishedItems {
				if item.Payload == "test1" {
					foundTest1 = true
				}
				if item.Payload == "test2" {
					foundTest2 = true
				}
			}
			assert.False(t, foundTest1)
			assert.True(t, foundTest2)
		},
	} {
		t.Run(tName, func(t *testing.T) {
			tFunc(t, pqGen())
		})
	}
}
