package db

import (
	"context"

	"github.com/zackarysantana/velocity/src/config"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Instance struct {
	Id primitive.ObjectID `bson:"_id,omitempty" json:"id"`

	Config config.Config `bson:"config,omitempty" json:"config"`
}

func (c *Connection) GetInstance(ctx context.Context, query interface{}) (*Instance, error) {
	var Instance Instance
	return &Instance, c.col("instances").FindOne(ctx, query).Decode(&Instance)
}

func (c *Connection) InsertInstance(ctx context.Context, config config.Config) (*Instance, error) {
	instance := Instance{
		Config: config,
	}

	r, err := c.col("instances").InsertOne(ctx, instance)
	if err != nil {
		return nil, err
	}
	instance.Id = r.InsertedID.(primitive.ObjectID)

	return &instance, nil
}

func (c *Connection) StartWorkflow(id primitive.ObjectID, workflow string) ([]*Job, error) {
	// add to workflow collection

	// add to job collection
	return nil, nil
}
