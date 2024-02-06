package db

import "go.mongodb.org/mongo-driver/mongo"

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

func (m *Mongo) GetUserByUsername(username string) (User, error) {
	panic("not implemented")
}
