package service

import (
	"context"
	"errors"
)

var (
	ErrEmptyQueue = errors.New("no items in queue")
)

// The provided type T is used as the Payload type and V is used as the ID type.
type PriorityQueue[ID any, Payload any] interface {
	// Push pushes an item into the given queue.
	Push(context.Context, string, ...PriorityQueueItem[Payload]) error
	// Pop pops an item from the given queue.
	Pop(context.Context, string) (PriorityQueuePoppedItem[ID, Payload], error)
}

type PriorityQueueItem[Payload any] struct {
	Priority int
	Payload  Payload
}

type PriorityQueuePoppedItem[ID any, Payload any] struct {
	Id      ID
	Payload Payload
}
