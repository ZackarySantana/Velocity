package service

import "errors"

var (
	ErrInvalidId = errors.New("id not found")
)

type IDCreator[T any] interface {
	Create() T
	Read(interface{}) (T, error)
	String(T) string
}
