package mongo

import (
	"github.com/zackarysantana/velocity/internal/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func NewIDCreator[T any]() service.IDCreator[T] {
	return &idCreator[T]{}
}

type idCreator[T any] struct{}

func (m *idCreator[T]) Create() T {
	var result any = primitive.NewObjectID()
	return result.(T)
}

func (m *idCreator[T]) Read(id interface{}) (T, error) {
	switch v := id.(type) {
	case primitive.ObjectID:
		return any(v).(T), nil
	case string:
		if objID, err := primitive.ObjectIDFromHex(v); err == nil {
			return any(objID).(T), nil
		}
	}
	return *new(T), service.ErrInvalidId
}

func (m *idCreator[T]) String(id T) string {
	if objID, ok := any(id).(primitive.ObjectID); ok {
		return objID.Hex()
	}
	return ""
}
