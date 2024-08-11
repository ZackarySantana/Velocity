package velocity

import (
	"encoding/json"
	"net/http"

	"github.com/zackarysantana/velocity/src/config"
)

type APIClient struct {
	*baseClient
}

func NewAPI(base string) *APIClient {
	return &APIClient{baseClient: newBaseClient(base)}
}

func (c *APIClient) Health() (*http.Response, error) {
	return c.do("GET", "/health", nil)
}

type APIStartRoutineRequest struct {
	Config  config.Config
	Routine string
}

type APIStartRoutineResponse struct {
	Id string `json:"id"`
}

// StartRoutine
func (c *APIClient) StartRoutine(config *config.Config, routine string) (*http.Response, *APIStartRoutineResponse, error) {
	resp, err := c.do("POST", "/routine/start", APIStartRoutineRequest{Config: *config, Routine: routine})
	if err != nil {
		return resp, nil, err
	}
	decodedResp := APIStartRoutineResponse{}
	defer resp.Body.Close()
	return resp, &decodedResp, json.NewDecoder(resp.Body).Decode(&decodedResp)
}
