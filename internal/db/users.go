package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id     primitive.ObjectID `bson:"_id,omitempty"`
	APIKey string             `bson:"api_key"`

	Email string `bson:"email"`
}

type Permissions struct {
	Id     primitive.ObjectID `bson:"_id,omitempty"`
	APIKey string             `bson:"api_key"`

	UserId string `bson:"user_id"`

	Admin bool `bson:"admin"`
}

func (u *User) CheckPassword(password string) bool {
	return true
}

func (c *Connection) GetUser(ctx context.Context, query interface{}) (*User, error) {
	var user User
	return &user, c.col("users").FindOne(ctx, query).Decode(&user)
}

func (c *Connection) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	return c.GetUser(ctx, bson.M{"email": email})
}

func (c *Connection) GetUserByAPIKey(ctx context.Context, apiKey string) (*User, error) {
	return c.GetUser(ctx, bson.M{"api_key": apiKey})
}

func (c *Connection) GetPermissions(ctx context.Context, query interface{}) (*Permissions, error) {
	var permissions Permissions
	return &permissions, c.col("permissions").FindOne(ctx, query).Decode(&permissions)
}

func (c *Connection) GetPermissionsByAPIKey(ctx context.Context, apiKey string) (*Permissions, error) {
	return c.GetPermissions(ctx, bson.M{"api_key": apiKey})
}
