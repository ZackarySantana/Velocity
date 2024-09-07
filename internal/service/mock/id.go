package mock

import (
	"github.com/google/uuid"
	"github.com/zackarysantana/velocity/internal/service"
)

func NewIDCreator[T any]() service.IDCreator[T] {
	return &mockIDCreator[T]{}
}

type mockIDCreator[T any] struct{}

func (m *mockIDCreator[T]) Create() T {
	var result any = uuid.New().String()
	return result.(T)
}

func (m *mockIDCreator[T]) Read(id interface{}) (T, error) {
	switch v := id.(type) {
	case string:
		return any(v).(T), nil
	}
	return *new(T), service.ErrInvalidId
}

func (m *mockIDCreator[T]) String(id T) string {
	return any(id).(string)
}
