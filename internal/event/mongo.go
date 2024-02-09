package event

import (
	"context"
	"time"

	"github.com/zackarysantana/velocity/internal/db"
	"github.com/zackarysantana/velocity/internal/event/meta"
	"go.mongodb.org/mongo-driver/mongo"
)

type Mongo struct {
	*mongo.Client

	db string
}

func NewMongo(client *mongo.Client, db string) EventSender {
	return &Mongo{
		Client: client,
		db:     db,
	}
}

func (m *Mongo) event() *mongo.Collection {
	return m.Database(m.db).Collection("events")
}

func (m *Mongo) sendEvent(ctx context.Context, event Event) error {
	event.TimeStamp = time.Now()
	_, err := m.event().InsertOne(ctx, event)
	return err
}

func (m *Mongo) SendIndexesAppliedEvent(ctx context.Context, user db.User) error {
	return m.sendEvent(ctx, Event{
		EventType: EventTypeIndexesApplied,
		Metadata:  meta.CreateApplyIndexes(user),
	})
}

func (m *Mongo) SendUserCreated(ctx context.Context, user db.User) error {
	return m.sendEvent(ctx, Event{
		EventType: EventTypeUserCreated,
		Metadata:  meta.CreateApplyIndexes(user),
	})
}
