package clients

import (
	"errors"
	"os"

	"github.com/zackarysantana/velocity/src/config"
)

type VelocityClientV1 struct {
	BaseURL string
}

func NewVelocityClientV1(baseURL string) *VelocityClientV1 {
	return &VelocityClientV1{
		BaseURL: baseURL,
	}
}

func NewVelocityClientV1FromEnv() (*VelocityClientV1, error) {
	baseURL := os.Getenv("VELOCITY_SERVER")
	if baseURL == "" {
		return nil, errors.New("VELOCITY_SERVER is not set")
	}

	return &VelocityClientV1{
		BaseURL: baseURL,
	}, nil
}

func NewVelocityClientV1FromConfig(c config.Config) (*VelocityClientV1, error) {
	if c.Config.Server == nil || *c.Config.Server == "" {
		return nil, errors.New("config does not have a server set")
	}

	return NewVelocityClientV1(*c.Config.Server), nil
}
