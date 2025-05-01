package service

import (
	"context"
	"io"
)

type ProcessQueue interface {
	io.Closer

	Write(ctx context.Context, topic string, messages ...[]byte) error
	Consume(ctx context.Context, topic string, consumerFunc func(item []byte) (processed bool, err error)) error
}
