package velocity

import (
	"encoding/json"
	"net/http"

	"github.com/zackarysantana/velocity/src/entities/test"
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

type AgentGetTestResponse[T any] struct {
	Test test.Test[T]
}

func (c *AgentClient) GetTest(id string) (*http.Response, *AgentGetTestResponse[any], error) {
	resp, err := c.do("GET", "/test/"+id, nil)
	if err != nil {
		return resp, nil, err
	}
	decodedResp := AgentGetTestResponse[any]{}
	defer resp.Body.Close()
	return resp, &decodedResp, json.NewDecoder(resp.Body).Decode(&decodedResp)
}
