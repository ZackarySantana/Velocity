package db

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`

	Username string `bson:"username" json:"username"`
	Password string `bson:"password" json:"password"`

	Email string `bson:"email" json:"email"`

	UserPermission UserPermission `bson:"permissions" json:"permissions"`
}

type UserPermission struct {
	SuperUser bool `bson:"super_user" json:"super_user"`
}
