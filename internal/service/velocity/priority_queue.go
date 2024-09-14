package velocity

import (
	"context"

	"github.com/samber/oops"
	"github.com/zackarysantana/velocity/internal/service"
	"github.com/zackarysantana/velocity/src/velocity"
)

func NewPriorityQueue[ID any, Payload any](client *velocity.AgentClient) service.PriorityQueue[ID, Payload] {
	return &priorityQueue[ID, Payload]{
		client: client,
	}
}

type priorityQueue[ID any, Payload any] struct {
	client *velocity.AgentClient
}

func (p *priorityQueue[ID, Payload]) Push(ctx context.Context, coll string, payloads ...service.PriorityQueueItem[Payload]) error {
	anyPayloads := make([]service.PriorityQueueItem[any], len(payloads))
	for i, payload := range payloads {
		anyPayloads[i] = service.PriorityQueueItem[any]{
			Priority: payload.Priority,
			Payload:  payload.Payload,
		}
	}
	_, err := p.client.Push(ctx, velocity.AgentPushRequest{
		Payloads: anyPayloads,
		Type:     coll,
	})
	return err
}

func (p *priorityQueue[ID, Payload]) Pop(ctx context.Context, coll string) (service.PriorityQueuePoppedItem[ID, Payload], error) {
	_, resp, err := p.client.Pop(ctx, velocity.AgentPopRequest{
		Type: coll,
	})
	if err != nil {
		return service.PriorityQueuePoppedItem[ID, Payload]{}, err
	}
	item := service.PriorityQueuePoppedItem[ID, Payload]{
		CreatedOn: resp.Popped.CreatedOn,
	}
	var ok bool
	item.Id, ok = resp.Popped.Id.(ID)
	if !ok {
		return service.PriorityQueuePoppedItem[ID, Payload]{}, oops.Errorf("could not convert %v to ID", resp.Popped.Id)
	}
	item.Payload, ok = resp.Popped.Payload.(Payload)
	if !ok {
		return service.PriorityQueuePoppedItem[ID, Payload]{}, oops.Errorf("could not convert %v to Payload", resp.Popped.Payload)
	}
	return item, nil
}

func (p *priorityQueue[ID, Payload]) MarkAsDone(ctx context.Context, coll string, id ID) error {
	_, err := p.client.MarkAsDone(ctx, velocity.AgentMarkAsDoneRequest{
		ID:   id,
		Type: coll,
	})
	return err
}

func (p *priorityQueue[ID, Payload]) UnfinishedItems(ctx context.Context, coll string) ([]service.PriorityQueueUnfinishedItem[ID, Payload], error) {
	_, resp, err := p.client.UnfinishedItems(ctx, velocity.AgentUnfinishedItemsRequest{
		Type: coll,
	})
	if err != nil {
		return nil, err
	}
	items := make([]service.PriorityQueueUnfinishedItem[ID, Payload], len(resp.Items))
	for i, item := range resp.Items {
		items[i] = service.PriorityQueueUnfinishedItem[ID, Payload]{
			CreatedOn: item.CreatedOn,
		}
		var ok bool
		items[i].Id, ok = item.Id.(ID)
		if !ok {
			return nil, oops.Errorf("could not convert %v to ID", item.Id)
		}
		items[i].Payload, ok = item.Payload.(Payload)
		if !ok {
			return nil, oops.Errorf("could not convert %v to Payload", item.Payload)
		}
	}
	return items, nil
}

func (p *priorityQueue[ID, Payload]) Close() error {
	return nil
}
