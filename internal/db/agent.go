package db

import "go.mongodb.org/mongo-driver/bson/primitive"

type Agent struct {
	Id primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`

	AgentSecret string `bson:"agent_secret" json:"agent_secret"`
}
