package mongo

import (
	"context"
	"fmt"
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

	CreatedOn time.Time
	StartedOn *time.Time
	EndedOn   *time.Time
}

func (p *priorityQueue[ID, Payload]) Push(ctx context.Context, coll string, payloads ...service.PriorityQueueItem[Payload]) error {
	items := make([]interface{}, len(payloads))
	for i, payload := range payloads {
		items[i] = queueItem[ID, Payload]{
			Id:        p.idCreator.Create(),
			Priority:  payload.Priority,
			Payload:   payload.Payload,
			CreatedOn: time.Now(),
			StartedOn: nil,
			EndedOn:   nil,
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
			"startedon": nil,
		},
		bson.M{
			"$set": bson.M{
				"startedon": time.Now(),
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
		Id:        item.Id,
		Payload:   item.Payload,
		CreatedOn: item.CreatedOn,
	}, err
}

func (p *priorityQueue[ID, Payload]) MarkAsDone(ctx context.Context, coll string, id ID) error {
	_, err := p.db.Database(p.dbName).Collection(coll).UpdateOne(ctx,
		bson.M{
			"_id": id,
		},
		bson.M{
			"$set": bson.M{
				"endedon": time.Now(),
			},
		},
	)
	return err
}

func (p *priorityQueue[ID, Payload]) UnfinishedItems(ctx context.Context, coll string) ([]service.PriorityQueueUnfinishedItem[ID, Payload], error) {
	items := []queueItem[ID, Payload]{}

	cursor, err := p.db.Database(p.dbName).Collection(coll).Find(ctx, bson.M{
		"startedon": bson.M{"$ne": nil},
		"endedon":   nil,
	})
	if err == mongo.ErrNoDocuments {
		return nil, service.ErrEmptyQueue
	}
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &items)
	if err != nil {
		return nil, err
	}

	unfinishedItems := make([]service.PriorityQueueUnfinishedItem[ID, Payload], len(items))
	for i, item := range items {
		if item.StartedOn == nil {
			return nil, fmt.Errorf("item %v has no start time", item.Id)
		}
		unfinishedItems[i] = service.PriorityQueueUnfinishedItem[ID, Payload]{
			Id:        item.Id,
			Payload:   item.Payload,
			CreatedOn: item.CreatedOn,
			StartedOn: *item.StartedOn,
		}
	}

	return unfinishedItems, err
}

func (p *priorityQueue[ID, Payload]) Close() error {
	return nil
}
