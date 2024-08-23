package mongo

import (
	"context"
	"time"

	"github.com/samber/oops"
	"github.com/zackarysantana/velocity/internal/service"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type queueItem[T any, V any] struct {
	Id       V `bson:"_id,omitempty"`
	Priority int

	Payload T

	StartTime *time.Time
	EndTime   *time.Time
	CreatedOn time.Time
}

func NewMongoPriorityQueue[T any, V any](db *mongo.Client, idCreator service.IdCreator[V], dbName string) service.PriorityQueue[T, V] {
	return &priorityQueue[T, V]{
		db:        db,
		dbName:    dbName,
		idCreator: idCreator,
	}
}

type priorityQueue[T any, V any] struct {
	db     *mongo.Client
	dbName string

	idCreator service.IdCreator[V]
}

func (p *priorityQueue[T, V]) Push(ctx context.Context, coll string, payloads ...service.PriorityQueueItem[T]) error {
	items := make([]interface{}, len(payloads))
	for i, payload := range payloads {
		items[i] = queueItem[T, V]{
			Id:        p.idCreator.Create(),
			Priority:  payload.Priority,
			Payload:   payload.Payload,
			StartTime: nil,
			EndTime:   nil,
			CreatedOn: time.Now(),
		}
		i++
	}
	_, err := p.db.Database(p.dbName).Collection(coll).InsertMany(ctx, items)
	return err
}

func (p *priorityQueue[T, V]) Pop(ctx context.Context, coll string) (service.PriorityQueuePoppedItem[T, V], error) {
	item := queueItem[T, V]{}
	err := p.db.Database(p.dbName).Collection(coll).FindOneAndUpdate(ctx,
		bson.M{
			"starttime": nil,
		},
		bson.M{
			"$set": bson.M{
				"starttime": time.Now(),
			},
		},
		options.FindOneAndUpdate().SetSort(
			bson.D{
				{Key: "priority", Value: -1},
				{Key: "createdon", Value: 1},
			},
		).SetReturnDocument(options.After),
	).Decode(&item)

	if err == mongo.ErrNoDocuments {
		return service.PriorityQueuePoppedItem[T, V]{}, oops.Errorf("no items in queue")
	}

	return service.PriorityQueuePoppedItem[T, V]{
		Id:      item.Id,
		Payload: item.Payload,
	}, err
}

func (p *priorityQueue[T, V]) CloseItem(ctx context.Context, coll string, payload T) error {
	_, err := p.db.Database(p.dbName).Collection(coll).UpdateOne(ctx,
		bson.M{
			"payload": payload,
		},
		bson.M{
			"$set": bson.M{
				"endtime": time.Now(),
			},
		},
	)
	return err
}

func (p *priorityQueue[T, V]) Close() error {
	return p.db.Disconnect(context.Background())
}
