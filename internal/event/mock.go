package event

import (
	"context"

	"github.com/zackarysantana/velocity/internal/db"
)

type Mock struct{}

func NewMock() EventSender {
	return &Mock{}
}

func (m *Mock) SendIndexesAppliedEvent(ctx context.Context, user db.User) error {
	return nil
}

func (m *Mock) SendUserCreated(ctx context.Context, user db.User) error {
	return nil
}
