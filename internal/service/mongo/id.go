package mongo

import (
	"github.com/zackarysantana/velocity/internal/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func NewObjectIDCreator[T any]() service.IDCreator[T] {
	return &mongoIDCreator[T]{}
}

type mongoIDCreator[T any] struct{}

func (m *mongoIDCreator[T]) Create() T {
	var result any = primitive.NewObjectID()
	return result.(T)
}

func (m *mongoIDCreator[T]) Read(id interface{}) (T, error) {
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

func (m *mongoIDCreator[T]) String(id T) string {
	if objID, ok := any(id).(primitive.ObjectID); ok {
		return objID.Hex()
	}
	return ""
}
