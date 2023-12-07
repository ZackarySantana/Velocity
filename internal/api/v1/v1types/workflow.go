package v1types

import "github.com/zackarysantana/velocity/src/config"

// POST /api/v1/workflows/start
type PostInstanceStartRequest struct {
	ProjectId string `json:"id"`

	Config   *config.Config `json:"config"`
	Workflow string         `json:"workflow"`
}
type PostInstanceStartResponse string

func NewPostInstanceStartRequest() interface{} {
	return &PostInstanceStartRequest{}
}
