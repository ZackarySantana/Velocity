package db

import (
	"context"

	"github.com/zackarysantana/velocity/src/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserInfo struct {
	Id    string `bson:"id,omitempty" json:"id,omitempty"`
	Email string `bson:"email,omitempty" json:"email,omitempty"`
}

type Instance struct {
	Id primitive.ObjectID `bson:"_id,omitempty" json:"id"`

	ProjectId primitive.ObjectID `bson:"project_id,omitempty" json:"project_id"`
	Config    config.Config      `bson:"config,omitempty" json:"config"`

	// Metadata
	UserInfo *UserInfo `bson:"user_info,omitempty" json:"user_info,omitempty"`
}

func (c *Connection) ApplyInstanceIndexes(ctx context.Context) error {
	_, err := c.col("instances").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.M{"project_id": 1},
	})
	return err
}

func (c *Connection) GetInstance(ctx context.Context, query interface{}) (*Instance, error) {
	var Instance Instance
	return &Instance, c.col("instances").FindOne(ctx, query).Decode(&Instance)
}

func (c *Connection) GetInstanceById(ctx context.Context, id primitive.ObjectID) (*Instance, error) {
	return c.GetInstance(ctx, bson.M{"_id": id})
}

func (c *Connection) InsertInstance(ctx context.Context, instance *Instance) (*Instance, error) {
	r, err := c.col("instances").InsertOne(ctx, instance)
	if err != nil {
		return nil, err
	}
	instance.Id = r.InsertedID.(primitive.ObjectID)

	return instance, nil
}
