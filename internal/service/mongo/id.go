package mongo

import (
	"github.com/zackarysantana/velocity/internal/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// NewObjectIDCreator creates a new IdCreator for MongoDB ObjectIDs.
func NewObjectIDCreator() service.IDCreator[primitive.ObjectID] {
	return &mongoIDCreator{}
}

type mongoIDCreator struct{}

func (m *mongoIDCreator) Create() primitive.ObjectID {
	return primitive.NewObjectID()
}

func (m *mongoIDCreator) Read(id interface{}) (primitive.ObjectID, error) {
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

func (m *mongoIDCreator) String(id primitive.ObjectID) string {
	return id.Hex()
}
