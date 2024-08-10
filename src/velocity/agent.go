package velocity

import (
	"net/http"
)

type AgentClient struct {
	*baseClient
}

func NewAgent(base string) *AgentClient {
	return &AgentClient{baseClient: newBaseClient(base + "/agent")}
}

func (c *AgentClient) Health() (*http.Response, error) {
	return c.do("GET", "/health", nil)
}
