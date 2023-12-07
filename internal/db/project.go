package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Project struct {
	Id primitive.ObjectID `bson:"_id,omitempty" json:"id"`

	Git struct {
		Owner      string `bson:"owner" json:"owner"`
		Repository string `bson:"repository" json:"repository"`
	} `bson:"git" json:"git"`
}

func (c *Connection) GetProject(ctx context.Context, query interface{}) (*Project, error) {
	var project Project
	return &project, c.col("projects").FindOne(ctx, query).Decode(&project)
}

func (c *Connection) InsertProject(ctx context.Context, project *Project) (*Project, error) {
	r, err := c.col("projects").InsertOne(ctx, project)
	if err != nil {
		return nil, err
	}
	project.Id = r.InsertedID.(primitive.ObjectID)

	return project, nil
}
