package service

import "context"

type ProcessQueue interface {
	Write(context.Context, string, ...[]byte) error
	Consume(context.Context, string, func([]byte) (bool, error)) error

	Close() error
}
