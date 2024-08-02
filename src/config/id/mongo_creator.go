package id

import "go.mongodb.org/mongo-driver/bson/primitive"

func NewMongoId() string {
	return primitive.NewObjectID().Hex()
}
