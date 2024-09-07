package mock

import (
	"context"
	"errors"
	"time"

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
	coll     string
	Payload  Payload

	CreatedOn time.Time
	StartedOn *time.Time
	EndedOn   *time.Time
}

func (p *priorityQueue[ID, Payload]) Push(ctx context.Context, coll string, payloads ...service.PriorityQueueItem[Payload]) error {
	for _, payload := range payloads {
		p.items = append(p.items, &queueItem[ID, Payload]{
			ID:        p.idCreator.Create(),
			Priority:  payload.Priority,
			coll:      coll,
			Payload:   payload.Payload,
			CreatedOn: time.Now(),
		})
	}

	return nil
}

func (p *priorityQueue[ID, Payload]) Pop(ctx context.Context, coll string) (service.PriorityQueuePoppedItem[ID, Payload], error) {
	highestPriority := -1

	for _, item := range p.items {
		if item.StartedOn != nil || item.coll != coll {
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
		if item.StartedOn != nil || item.coll != coll || item.Priority != highestPriority {
			continue
		}

		t := time.Now()
		item.StartedOn = &t

		return service.PriorityQueuePoppedItem[ID, Payload]{
			Id:        item.ID,
			Payload:   item.Payload,
			CreatedOn: item.CreatedOn,
		}, nil
	}

	return service.PriorityQueuePoppedItem[ID, Payload]{}, service.ErrEmptyQueue
}

func (p *priorityQueue[ID, Payload]) MarkAsDone(ctx context.Context, coll string, id ID) error {
	for _, item := range p.items {
		if item.ID == id {
			t := time.Now()
			item.EndedOn = &t
			return nil
		}
	}
	return errors.New("item not found")
}

func (p *priorityQueue[ID, Payload]) UnfinishedItems(ctx context.Context, coll string) ([]service.PriorityQueueUnfinishedItem[ID, Payload], error) {
	unfinsihedItems := []service.PriorityQueueUnfinishedItem[ID, Payload]{}

	for _, item := range p.items {
		if item.coll == coll && item.StartedOn != nil && item.EndedOn == nil {
			unfinsihedItems = append(unfinsihedItems, service.PriorityQueueUnfinishedItem[ID, Payload]{
				Id:      item.ID,
				Payload: item.Payload,
			})
		}
	}

	return unfinsihedItems, nil
}

func (p *priorityQueue[ID, Payload]) Close() error {
	return nil
}
