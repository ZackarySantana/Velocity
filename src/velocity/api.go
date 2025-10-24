package velocity

import (
	"context"
	"fmt"
	"net/http"

	"github.com/zackarysantana/velocity/src/config"
	"github.com/zackarysantana/velocity/src/entities/routine"
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
	Config  config.Config `json:"config"`
	Routine string        `json:"routine"`
}

func (a *APIStartRoutineRequest) Validate() error {
	if err := a.Config.Validate(); err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}
	if a.Config.Routines.GetRoutine(a.Routine) == nil {
		return fmt.Errorf("routine '%s' not found", a.Routine)
	}
	return nil
}

type APIStartRoutineResponse struct {
	Id interface{} `json:"id"`
}

func (c *APIClient) StartRoutine(ctx context.Context, config *config.Config, routine string) (*http.Response, *APIStartRoutineResponse, error) {
	decodedResp := APIStartRoutineResponse{}
	resp, err := c.doAndDecode(ctx, "POST", "/routine/start", APIStartRoutineRequest{Config: *config, Routine: routine}, &decodedResp)
	return resp, &decodedResp, err
}

type APIGetRoutineResponse[T any] struct {
	Routine routine.Routine[T] `json:"routine"`
}

func (c *APIClient) GetRoutine(ctx context.Context, id interface{}) (*http.Response, *APIGetRoutineResponse[any], error) {
	decodedResp := APIGetRoutineResponse[any]{}
	resp, err := c.doAndDecode(ctx, "GET", fmt.Sprintf("/routine/%s", id), nil, &decodedResp)
	return resp, &decodedResp, err
}
