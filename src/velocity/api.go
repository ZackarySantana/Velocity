package velocity

import (
	"context"
	"net/http"

	"github.com/zackarysantana/velocity/src/config"
)

type APIClient struct {
	*baseClient
}

func NewAPIClient(base string) *APIClient {
	return &APIClient{baseClient: newBaseClient(base)}
}

func (c *APIClient) Health(ctx context.Context) (*http.Response, error) {
	return c.do(ctx, "GET", "/health", nil)
}

type APIStartRoutineRequest struct {
	Config  config.Config
	Routine string
}

type APIStartRoutineResponse struct {
	Id interface{} `json:"id"`
}

func (c *APIClient) StartRoutine(ctx context.Context, config *config.Config, routine string) (*http.Response, *APIStartRoutineResponse, error) {
	decodedResp := APIStartRoutineResponse{}
	resp, err := c.doAndDecode(ctx, "POST", "/routine/start", APIStartRoutineRequest{Config: *config, Routine: routine}, &decodedResp)
	return resp, &decodedResp, err
}
