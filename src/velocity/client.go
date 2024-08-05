package velocity

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/samber/oops"
	"github.com/zackarysantana/velocity/src/config"
)

type Client struct {
	client *http.Client
	base   string
}

func New(base string) *Client {
	return &Client{
		client: &http.Client{},
		base:   base,
	}
}

func (c *Client) Health() (*http.Response, error) {
	return c.do("GET", "/health", nil)
}

type StartRoutineRequst struct {
	Config  config.Config
	Routine string
}

type StartRoutineResponse struct {
	Id string `json:"id"`
}

// StartRoutine
func (c *Client) StartRoutine(config *config.Config, routine string) (*http.Response, *StartRoutineResponse, error) {
	resp, err := c.do("POST", "/routine/start", StartRoutineRequst{Config: *config, Routine: routine})
	if err != nil {
		return resp, nil, err
	}
	decodedResp := StartRoutineResponse{}
	defer resp.Body.Close()
	return resp, &decodedResp, json.NewDecoder(resp.Body).Decode(&decodedResp)
}

func (c *Client) do(method, path string, payload interface{}) (*http.Response, error) {
	var body *bytes.Buffer
	if payload != nil {
		body = &bytes.Buffer{}
		err := json.NewEncoder(body).Encode(payload)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, c.base+path, body)
	if err != nil {
		return nil, err
	}
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return resp, err
	}
	if resp.StatusCode != http.StatusOK {
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return resp, oops.Code("decoding").Wrap(err)
		}
		return resp, oops.Code("status").Wrap(oops.With("status", resp.StatusCode).Errorf("unexpected error: %s", string(respBody)))
	}
	return resp, nil
}
