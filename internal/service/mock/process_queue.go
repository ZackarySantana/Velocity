package mock

import (
	"context"
	"errors"
	"sync"

	"github.com/zackarysantana/velocity/internal/service"
)

var _ service.ProcessQueue = (*processQueue)(nil)

func NewProcessQueue() *processQueue {
	return &processQueue{
		items: make(map[string][][]byte),
		cond:  sync.NewCond(&sync.Mutex{}),
	}
}

type processQueue struct {
	cond   *sync.Cond
	items  map[string][][]byte
	closed bool
}

func (p *processQueue) Write(ctx context.Context, topic string, messages ...[]byte) error {
	p.cond.L.Lock()
	defer p.cond.L.Unlock()

	if p.closed {
		return errors.New("queue is closed")
	}

	p.items[topic] = append(p.items[topic], messages...)
	p.cond.Broadcast() // wake consumers
	return nil
}

func (p *processQueue) Consume(ctx context.Context, topic string, consumerFunc func(message []byte) (processed bool, err error)) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			p.cond.L.Lock()

			// Wait for new items or close
			for !p.closed && len(p.items[topic]) == 0 {
				p.cond.Wait()
			}

			if p.closed {
				p.cond.L.Unlock()
				return nil
			}

			// Dequeue the item
			item := p.items[topic][0]
			p.items[topic] = p.items[topic][1:]
			p.cond.L.Unlock()

			processed, err := consumerFunc(item)
			if err != nil {
				return err
			}

			// if not processed, put it back at the top of the queue
			if !processed {
				p.cond.L.Lock()
				p.items[topic] = append([][]byte{item}, p.items[topic]...)
				p.cond.L.Unlock()
				continue
			}
		}
	}
}

func (p *processQueue) Close() error {
	p.cond.L.Lock()
	defer p.cond.L.Unlock()

	if p.closed {
		return errors.New("queue already closed")
	}

	p.closed = true
	p.cond.Broadcast()
	return nil
}
