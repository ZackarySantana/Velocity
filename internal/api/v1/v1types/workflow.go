package v1types

import "github.com/zackarysantana/velocity/src/config"

// POST /api/v1/workflows/start
type PostWorkflowsStartRequest struct {
	Config config.YAMLWorkflow `json:"config"`
}
type PostWorkflowsStartResponse string
