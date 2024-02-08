package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Mongo struct {
	*mongo.Client

	db string
}

func NewMongo(client *mongo.Client, db string) Database {
	return &Mongo{
		Client: client,
		db:     db,
	}
}

func (m *Mongo) user() *mongo.Collection {
	return m.Database(m.db).Collection("users")
}

func (m *Mongo) agent() *mongo.Collection {
	return m.Database(m.db).Collection("agents")
}

func (m *Mongo) ApplyIndexes(ctx context.Context) error {
	err := m.ApplyUserIndexes(ctx)
	if err != nil {
		return err
	}
	err = m.ApplyAgentIndexes(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (m *Mongo) ApplyUserIndexes(ctx context.Context) error {
	_, err := m.user().Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "username", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return err
	}
	_, err = m.user().Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "email", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return err
	}
	return nil
}

func (m *Mongo) ApplyAgentIndexes(ctx context.Context) error {
	_, err := m.agent().Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "agent_secret", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return err
	}
	return nil
}

func (m *Mongo) GetUserByUsername(ctx context.Context, username string) (User, error) {
	var user User
	err := m.user().FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return user, ErrNoEntity
	}
	return user, err
}

func (m *Mongo) CreateUser(ctx context.Context, user User) (User, error) {
	id, err := m.user().InsertOne(ctx, user)
	if err != nil {
		return user, err
	}
	user.Id = id.InsertedID.(primitive.ObjectID)
	return user, nil
}

func (m *Mongo) GetAgentBySecret(ctx context.Context, agentSecret string) (Agent, error) {
	var agent Agent
	err := m.agent().FindOne(ctx, bson.M{"agent_secret": agentSecret}).Decode(&agent)
	return agent, err
}
