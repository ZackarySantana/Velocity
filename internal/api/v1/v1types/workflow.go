package v1types

import "github.com/zackarysantana/velocity/src/config"

type PostWorkflowsStartRequest struct {
	Config config.YAMLWorkflow `json:"config"`
}

// TBA
type PostWorkflowsStartResponse string
