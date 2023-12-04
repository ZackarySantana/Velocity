package v1types

import "github.com/zackarysantana/velocity/src/config"

type PostWorkflowRequest struct {
	Config config.YAMLWorkflow `json:"config"`
}

// TBA
type PostWorkflowResponse string
