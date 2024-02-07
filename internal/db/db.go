package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id primitive.ObjectID `bson:"_id,omitempty"`

	Username string `bson:"username"`
	Password string `bson:"password"`

	Email string `bson:"email"`
}

type Agent struct {
	Id primitive.ObjectID `bson:"_id,omitempty"`

	AgentSecret string `bson:"agent_secret"`
}

type Database interface {
	GetUserByUsername(ctx context.Context, username string) (User, error)
	GetAgentBySecret(ctx context.Context, agentSecret string) (Agent, error)
}
