package mongo

import (
	"github.com/zackarysantana/velocity/internal/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func NewMongoIdCreator() service.IdCreator[primitive.ObjectID] {
	return &mongoIdCreator{}
}

type mongoIdCreator struct{}

func (m *mongoIdCreator) Create() primitive.ObjectID {
	return primitive.NewObjectID()
}

func (m *mongoIdCreator) Read(id interface{}) (primitive.ObjectID, error) {
	objectId, ok := id.(primitive.ObjectID)
	if ok {
		return objectId, nil
	}
	str, ok := id.(string)
	if ok {
		if id, err := primitive.ObjectIDFromHex(str); err == nil {
			return id, nil
		}
	}
	return primitive.ObjectID{}, service.ErrInvalidId
}

func (m *mongoIdCreator) String(id primitive.ObjectID) string {
	return id.Hex()
}
