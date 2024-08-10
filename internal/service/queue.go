package service

import "context"

type ProcessQueue interface {
	Write(context.Context, string, ...[]byte) error
	Consume(context.Context, string, func([]byte) error) error

	Close() error
}
