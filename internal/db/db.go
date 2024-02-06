package db

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id primitive.ObjectID `bson:"_id,omitempty"`

	Username string `bson:"username"`
	Password string `bson:"password"`

	Email string `bson:"email"`
}

type Database interface {
	GetUserByUsername(username string) (User, error)
}
