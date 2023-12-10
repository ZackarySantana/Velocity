package clients

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/google/go-querystring/query"
	"github.com/zackarysantana/velocity/internal/api/v1/v1types"
	"github.com/zackarysantana/velocity/src/config"
)

type VelocityClientV1 struct {
	BaseURL string

	ApiKey string
}

func NewVelocityClientV1(baseURL string) *VelocityClientV1 {
	return &VelocityClientV1{
		BaseURL: strings.TrimSuffix(baseURL, "/"),
		ApiKey:  "YOUR_API_KEY",
	}
}

func NewVelocityClientV1WithAPIKey(baseURL, apiKey string) *VelocityClientV1 {
	return &VelocityClientV1{
		BaseURL: strings.TrimSuffix(baseURL, "/"),
		ApiKey:  apiKey,
	}
}

func NewVelocityClientV1FromEnv() (*VelocityClientV1, error) {
	baseURL := os.Getenv("VELOCITY_SERVER")
	if baseURL == "" {
		return nil, errors.New("VELOCITY_SERVER is not set")
	}

	return NewVelocityClientV1(baseURL), nil
}

func NewVelocityClientV1FromConfig(c *config.Config) (*VelocityClientV1, error) {
	if c.Config.Server == nil || *c.Config.Server == "" {
		return nil, errors.New("config does not have a server set")
	}

	return NewVelocityClientV1(*c.Config.Server), nil
}

func (v *VelocityClientV1) PostFirstTimeRegister(body v1types.PostFirstTimeRegisterRequest) (*v1types.PostFirstTimeRegisterResponse, error) {
	var data v1types.PostFirstTimeRegisterResponse
	return &data, v.post("/first_time_register", body, &data)
}

func (v *VelocityClientV1) PostInstanceStart(body v1types.PostInstanceStartRequest) (*v1types.PostInstanceStartResponse, error) {
	var data v1types.PostInstanceStartResponse
	return &data, v.post("/instances/start", body, &data)
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

	path := v.makeRoute(route)
	req, err := http.NewRequest("POST", path, bytes.NewBuffer(out))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+v.ApiKey)

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		b, _ := ioutil.ReadAll(response.Body)
		return errors.New(path + " - " + response.Status + ": " + string(b))
	}

	ct := response.Header.Get("Content-Type")
	if !strings.Contains(ct, "application/json") {
		return errors.New(path + " - " + ct + ": invalid content type")
	}

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
