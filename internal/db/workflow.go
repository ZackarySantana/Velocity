package db

import "go.mongodb.org/mongo-driver/bson/primitive"

type Workflow struct {
	Id primitive.ObjectID `bson:"_id,omitempty"`

	Name   string          `bson:"name"`
	Groups []WorkflowGroup `bson:"groups"`
}

type WorkflowGroup struct {
	Name string `bson:"name"`

	RuntimeIDs []primitive.ObjectID `bson:"runtime_ids"`
	TestIDs    []primitive.ObjectID `bson:"test_ids"`
}
