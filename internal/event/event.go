package event

import (
	"context"
	"time"

	"github.com/zackarysantana/velocity/internal/db"
)

// TODO: Switch from the db.User to a custom type or interface that is more generic and flexible.
// Also helps keep the event package from being too tightly coupled to the db package.
type EventSender interface {
	SendIndexesAppliedEvent(ctx context.Context, user db.User) error
	SendUserCreated(ctx context.Context, user db.User) error
}

type EventType string

const (
	EventTypeUserCreated    EventType = "user_created"
	EventTypeIndexesApplied EventType = "indexes_applied"
)

type Event struct {
	EventType EventType `json:"event_type" bson:"event_type"`

	TimeStamp time.Time         `json:"timestamp" bson:"timestamp"`
	Metadata  map[string]string `json:"metadata" bson:"metadata"`
}
