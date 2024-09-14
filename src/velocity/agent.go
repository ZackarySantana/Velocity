package velocity

import (
	"context"
	"net/http"

	"github.com/zackarysantana/velocity/internal/service"
	"github.com/zackarysantana/velocity/src/entities/test"
)

type AgentClient struct {
	*baseClient
}

func NewAgentClient(base string) *AgentClient {
	return &AgentClient{baseClient: newBaseClient(base + "/agent")}
}

func (c *AgentClient) Health(ctx context.Context) (*http.Response, error) {
	return c.do(ctx, "GET", "/health", nil)
}

type AgentGetTestResponse[T any] struct {
	Test test.Test[T]
}

func (c *AgentClient) GetTest(ctx context.Context, id string) (*http.Response, *AgentGetTestResponse[any], error) {
	decodedResp := AgentGetTestResponse[any]{}
	resp, err := c.doAndDecode(ctx, "GET", "/test/"+id, nil, &decodedResp)
	return resp, &decodedResp, err
}

type AgentPushRequest struct {
	Type     string                           `json:"type"`
	Payloads []service.PriorityQueueItem[any] `json:"payloads"`
}

func (c *AgentClient) Push(ctx context.Context, req AgentPushRequest) (*http.Response, error) {
	return c.do(ctx, "POST", "/priority_queue/push", req)
}

type AgentPopRequest struct {
	Type string `json:"type"`
}

type AgentPopResponse struct {
	Popped service.PriorityQueuePoppedItem[any, any] `json:"payloads"`
}

func (c *AgentClient) Pop(ctx context.Context, req AgentPopRequest) (*http.Response, *AgentPopResponse, error) {
	decodedResp := AgentPopResponse{}
	resp, err := c.doAndDecode(ctx, "POST", "/priority_queue/pop", req, &decodedResp)
	return resp, &decodedResp, err
}

type AgentMarkAsDoneRequest struct {
	ID   any    `json:"id"`
	Type string `json:"type"`
}

func (c *AgentClient) MarkAsDone(ctx context.Context, req AgentMarkAsDoneRequest) (*http.Response, error) {
	return c.do(ctx, "POST", "/priority_queue/done", req)
}

type AgentUnfinishedItemsRequest struct {
	Type string `json:"type"`
}

type AgentUnfinishedItemsResponse struct {
	Items []service.PriorityQueueUnfinishedItem[any, any] `json:"items"`
}

func (c *AgentClient) UnfinishedItems(ctx context.Context, req AgentUnfinishedItemsRequest) (*http.Response, *AgentUnfinishedItemsResponse, error) {
	decodedResp := AgentUnfinishedItemsResponse{}
	resp, err := c.doAndDecode(ctx, "POST", "/priority_queue/unfinished", req, &decodedResp)
	return resp, &decodedResp, err
}
