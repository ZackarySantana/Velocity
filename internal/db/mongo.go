package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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

func (m *Mongo) GetUserByUsername(ctx context.Context, username string) (User, error) {
	var user User
	err := m.user().FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return user, ErrNoEntity
	}
	return user, err
}

func (m *Mongo) GetAgentBySecret(ctx context.Context, agentSecret string) (Agent, error) {
	var agent Agent
	err := m.agent().FindOne(ctx, bson.M{"agent_secret": agentSecret}).Decode(&agent)
	return agent, err
}
