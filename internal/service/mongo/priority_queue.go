package mongo

import (
	"context"
	"time"

	"github.com/zackarysantana/velocity/internal/service"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// NewPriorityQueue creates a new PriorityQueue that uses MongoDB as the backend.
// The provided type T is used as the Payload type and V is used as the ID type.
func NewPriorityQueue[ID any, Payload any](db *mongo.Client, idCreator service.IDCreator[ID], dbName string) service.PriorityQueue[ID, Payload] {
	return &priorityQueue[ID, Payload]{
		db:        db,
		dbName:    dbName,
		idCreator: idCreator,
	}
}

type priorityQueue[ID any, Payload any] struct {
	db     *mongo.Client
	dbName string

	idCreator service.IDCreator[ID]
}

type queueItem[ID any, Payload any] struct {
	Id       ID `bson:"_id,omitempty"`
	Priority int

	Payload Payload

	StartTime *time.Time
	EndTime   *time.Time
	CreatedOn time.Time
}

func (p *priorityQueue[ID, Payload]) Push(ctx context.Context, coll string, payloads ...service.PriorityQueueItem[Payload]) error {
	items := make([]interface{}, len(payloads))
	for i, payload := range payloads {
		items[i] = queueItem[ID, Payload]{
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

func (p *priorityQueue[ID, Payload]) Pop(ctx context.Context, coll string) (service.PriorityQueuePoppedItem[ID, Payload], error) {
	item := queueItem[ID, Payload]{}
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
		return service.PriorityQueuePoppedItem[ID, Payload]{}, service.ErrEmptyQueue
	}

	return service.PriorityQueuePoppedItem[ID, Payload]{
		Id:      item.Id,
		Payload: item.Payload,
	}, err
}

func (p *priorityQueue[ID, Payload]) MarkAsDone(ctx context.Context, coll string, id ID) error {
	_, err := p.db.Database(p.dbName).Collection(coll).UpdateOne(ctx,
		bson.M{
			"_id": id,
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
	return nil
}
