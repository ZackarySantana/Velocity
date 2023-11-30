package clients

import (
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/zackarysantana/velocity/internal/db"
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

func (v *VelocityClientV1) Login(username, password string) (*db.User, error) {
	response, err := http.Get(v.makeRoute("/login"))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// Parse body
	return nil, nil
}

func (v *VelocityClientV1) makeRoute(route string) string {
	return v.BaseURL + "/api/v1/" + strings.TrimPrefix(route, "/")
}
