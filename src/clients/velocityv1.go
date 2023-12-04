package clients

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"

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

func (v *VelocityClientV1) PostWorkflow(body v1types.PostWorkflowRequest) (*v1types.PostWorkflowResponse, error) {
	var data v1types.PostWorkflowResponse
	return &data, v.post("/workflows", body, &data)
}

func (v *VelocityClientV1) PostFirstTimeRegister(email string) (*v1types.PostRegisterUserResponse, error) {
	var data v1types.PostJobResultRequest
	return &data, v.post("/workflows", body, &data)
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
	return v.BaseURL + "/api/v1/" + strings.TrimPrefix(route, "/")
}
