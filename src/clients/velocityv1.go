package clients

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/google/go-querystring/query"
	"github.com/zackarysantana/velocity/internal/api/v1/v1types"
	"github.com/zackarysantana/velocity/src/config"
)

type VelocityClientV1 struct {
	BaseURL string

	SessionToken *string
}

func NewVelocityClientV1(baseURL string) *VelocityClientV1 {
	return &VelocityClientV1{
		BaseURL: strings.TrimSuffix(baseURL, "/"),
	}
}

func NewVelocityClientV1WithSessionToken(baseURL, sessionToken string) *VelocityClientV1 {
	return &VelocityClientV1{
		BaseURL:      strings.TrimSuffix(baseURL, "/"),
		SessionToken: &sessionToken,
	}
}

func NewVelocityClientV1FromEnv() (*VelocityClientV1, error) {
	baseURL := os.Getenv("VELOCITY_SERVER")
	if baseURL == "" {
		return nil, errors.New("VELOCITY_SERVER is not set")
	}

	return NewVelocityClientV1(baseURL), nil
}

func NewVelocityClientV1FromConfig(c config.Config) (*VelocityClientV1, error) {
	if c.Config.Server == nil || *c.Config.Server == "" {
		return nil, errors.New("config does not have a server set")
	}

	return NewVelocityClientV1(*c.Config.Server), nil
}

func (v *VelocityClientV1) PostFirstTimeRegister(body v1types.PostFirstTimeRegisterRequest) (*v1types.PostFirstTimeRegisterResponse, error) {
	var data v1types.PostFirstTimeRegisterResponse
	return &data, v.post("/first_time_register", body, &data)
}

func (v *VelocityClientV1) PostWorkflowsStart(body v1types.PostWorkflowsStartRequest) (*v1types.PostWorkflowsStartResponse, error) {
	var data v1types.PostWorkflowsStartResponse
	return &data, v.post("/workflows", body, &data)
}

func (v *VelocityClientV1) PostJobsDequeue(body v1types.PostJobsDequeueRequest, opts v1types.PostJobsDequeueQueryParams) (*v1types.PostJobsDequeueResponse, error) {
	q, err := parseQueryParams(opts)
	if err != nil {
		return nil, err
	}
	var data v1types.PostJobsDequeueResponse
	return &data, v.post("/jobs/dequeue"+q, body, &data)
}

func (v *VelocityClientV1) PostJobsResults(body v1types.PostJobResultRequest) (*v1types.PostJobResultResponse, error) {
	var data v1types.PostJobResultResponse
	return &data, v.post("/jobs/result", body, &data)
}

func (v *VelocityClientV1) post(route string, body, resp interface{}) error {
	out, err := json.Marshal(body)
	if err != nil {
		return err
	}

	response, err := http.Post(v.makeRoute(route), "application/json", bytes.NewBuffer(out))
	if err != nil {
		return err
	}
	defer response.Body.Close()

	return json.NewDecoder(response.Body).Decode(&resp)
}

func (v *VelocityClientV1) makeRoute(route string) string {
	return v.BaseURL + "/api/v1" + route
}

func parseQueryParams(opts interface{}) (string, error) {
	q, err := query.Values(opts)
	if err != nil {
		return "", err
	}
	qp := q.Encode()
	if qp == "" {
		return "", nil
	}
	return "?" + q.Encode(), nil
}
