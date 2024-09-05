package service

import (
	"context"
)

type PriorityQueue[T any, V any] interface {
	// Push pushes an item into the given queue.
	Push(context.Context, string, ...PriorityQueueItem[T]) error
	// Pop pops an item from the given queue.
	Pop(context.Context, string) (PriorityQueuePoppedItem[T, V], error)
}

type PriorityQueueItem[T any] struct {
	Priority int
	Payload  T
}

type PriorityQueuePoppedItem[T any, V any] struct {
	Id      V
	Payload T
}