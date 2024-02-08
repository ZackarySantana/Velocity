package db

import "go.mongodb.org/mongo-driver/bson/primitive"

type Test struct {
	Id primitive.ObjectID `bson:"_id,omitempty"`

	Name      string             `bson:"name"`
	RuntimeID primitive.ObjectID `bson:"runtime_id"`
}
