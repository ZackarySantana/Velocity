package mock

import (
	"context"
	"errors"

	"github.com/zackarysantana/velocity/internal/service"
)

func NewPriorityQueue[ID comparable, Payload any](idCreator service.IDCreator[ID]) service.PriorityQueue[ID, Payload] {
	return &priorityQueue[ID, Payload]{
		idCreator: idCreator,
		items:     []*queueItem[ID, Payload]{},
	}
}

type priorityQueue[ID comparable, Payload any] struct {
	idCreator service.IDCreator[ID]
	items     []*queueItem[ID, Payload]
}

type queueItem[ID comparable, Payload any] struct {
	ID       ID
	Priority int
	Payload  Payload
	Started  bool
}

func (p *priorityQueue[ID, Payload]) Push(ctx context.Context, coll string, payloads ...service.PriorityQueueItem[Payload]) error {
	for _, payload := range payloads {
		p.items = append(p.items, &queueItem[ID, Payload]{
			ID:       p.idCreator.Create(),
			Priority: payload.Priority,
			Payload:  payload.Payload,
			Started:  false,
		})
	}

	return nil
}

func (p *priorityQueue[ID, Payload]) Pop(ctx context.Context, coll string) (service.PriorityQueuePoppedItem[ID, Payload], error) {
	highestPriority := -1

	for _, item := range p.items {
		if item.Started {
			continue
		}
		if item.Priority > highestPriority {
			highestPriority = item.Priority
		}
	}

	if highestPriority == -1 {
		return service.PriorityQueuePoppedItem[ID, Payload]{}, service.ErrEmptyQueue
	}

	for _, item := range p.items {
		if item.Priority != highestPriority || item.Started {
			continue
		}

		item.Started = true

		return service.PriorityQueuePoppedItem[ID, Payload]{
			Id:      item.ID,
			Payload: item.Payload,
		}, nil
	}

	return service.PriorityQueuePoppedItem[ID, Payload]{}, service.ErrEmptyQueue
}

func (p *priorityQueue[ID, Payload]) MarkAsDone(ctx context.Context, coll string, id ID) error {
	for _, item := range p.items {
		if item.ID == id {
			item.Started = false
			break
		}
	}
	return errors.New("item not found")
}

func (p *priorityQueue[ID, Payload]) Close() error {
	return nil
}
