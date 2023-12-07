package v1types

import (
	"github.com/zackarysantana/velocity/internal/db"
	"github.com/zackarysantana/velocity/src/config"
)

// POST /api/v1/workflows/start
type PostInstanceStartRequest struct {
	ProjectId string `json:"project_id"`

	Config   *config.Config `json:"config"`
	Workflow string         `json:"workflow"`
}
type PostInstanceStartResponse struct {
	InstanceId string `json:"instance_id"`

	Jobs []*db.Job
}

func NewPostInstanceStartRequest() interface{} {
	return &PostInstanceStartRequest{}
}
