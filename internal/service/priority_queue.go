package service

import (
	"context"
	"errors"
	"time"
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
	// MarkAsDone marks an item as done.
	MarkAsDone(context.Context, string, ID) error
	// UnfinishedItems returns all unfinished items.
	UnfinishedItems(context.Context, string) ([]PriorityQueueUnfinishedItem[ID, Payload], error)
	// Close closes the queue.
	Close() error
}

type PriorityQueueItem[Payload any] struct {
	Priority int
	Payload  Payload
}

type PriorityQueuePoppedItem[ID any, Payload any] struct {
	Id        ID
	Payload   Payload
	CreatedOn time.Time
}

type PriorityQueueUnfinishedItem[ID any, Payload any] struct {
	Id        ID
	Payload   Payload
	CreatedOn time.Time
	StartedOn time.Time
}
