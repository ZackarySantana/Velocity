package v1types

import "github.com/zackarysantana/velocity/src/config"

// POST /api/v1/workflows/start
type PostWorkflowsStartRequest struct {
	Config   *config.Config `json:"config"`
	Workflow string         `json:"workflow"`
}
type PostWorkflowsStartResponse string

func NewPostWorkflowsStartRequest() interface{} {
	return &PostWorkflowsStartRequest{}
}
