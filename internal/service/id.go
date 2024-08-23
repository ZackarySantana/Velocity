package service

import "errors"

var (
	ErrInvalidId = errors.New("id not found")
)

type IdCreator[T any] interface {
	Create() T
	Read(interface{}) (T, error)
	String(T) string
}
