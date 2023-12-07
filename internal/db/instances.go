package db

import (
	"context"

	"github.com/zackarysantana/velocity/src/config"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Instance struct {
	Id primitive.ObjectID `bson:"_id,omitempty" json:"id"`

	ProjectId primitive.ObjectID `bson:"project_id,omitempty" json:"project_id"`

	Config config.Config `bson:"config,omitempty" json:"config"`
}

func (c *Connection) GetInstance(ctx context.Context, query interface{}) (*Instance, error) {
	var Instance Instance
	return &Instance, c.col("instances").FindOne(ctx, query).Decode(&Instance)
}

func (c *Connection) InsertInstance(ctx context.Context, instance *Instance) (*Instance, error) {
	r, err := c.col("instances").InsertOne(ctx, instance)
	if err != nil {
		return nil, err
	}
	instance.Id = r.InsertedID.(primitive.ObjectID)

	return instance, nil
}
