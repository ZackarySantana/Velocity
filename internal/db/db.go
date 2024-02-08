package db

import (
	"context"
	"errors"
)

var ErrNoEntity = errors.New("no entity found")

type Database interface {
	// User
	GetUserByUsername(ctx context.Context, username string) (User, error)
	CreateUser(ctx context.Context, user User) (User, error)

	GetAgentBySecret(ctx context.Context, agentSecret string) (Agent, error)
}
