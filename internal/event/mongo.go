package event

import (
	"context"
	"time"

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

func (m *Mongo) SendEvent(ctx context.Context, event Event) error {
	event.TimeStamp = time.Now()
	_, err := m.event().InsertOne(ctx, event)
	return err
}
