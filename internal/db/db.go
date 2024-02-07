package db

import (
	"context"
	"errors"
)

var ErrNoEntity = errors.New("no entity found")

type Database interface {
	GetUserByUsername(ctx context.Context, username string) (User, error)
	GetAgentBySecret(ctx context.Context, agentSecret string) (Agent, error)
}
