package event

import (
	"context"
	"time"
)

type EventSender interface {
	SendEvent(ctx context.Context, event Event) error
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
