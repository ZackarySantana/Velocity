package mongo

import (
	"github.com/zackarysantana/velocity/internal/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func NewMongoIdCreator() service.IdCreator {
	return newMongoId
}

func newMongoId() string {
	return primitive.NewObjectID().Hex()
}
