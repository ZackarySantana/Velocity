package v1types

import (
	"github.com/zackarysantana/velocity/internal/db"
	"github.com/zackarysantana/velocity/src/config"
)

// GET /api/v1/instances
type GetInstanceRequest struct{}
type GetInstanceResponse struct {
	Instance db.Instance `json:"instance"`

	Jobs []db.Job `json:"jobs"`
}

// POST /api/v1/instances/start
type PostInstanceStartRequest struct {
	ProjectId   *string `json:"project_id"`
	ProjectName *string `json:"project_name"`

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
