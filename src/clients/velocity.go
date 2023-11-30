package clients

import (
	"errors"
	"os"
)

type Version string

const (
	v1 Version = "v1"
)

var (
	versions = []Version{v1}
)

type VelocityClient struct {
	BaseURL string

	APIVersion Version
}

func NewVelocityClientV1(baseURL string) *VelocityClient {
	return &VelocityClient{
		BaseURL:    baseURL,
		APIVersion: v1,
	}
}

func NewVelocityClientFromEnv() (*VelocityClient, error) {
	baseURL := os.Getenv("VELOCITY_SERVER")
	if baseURL == "" {
		return nil, errors.New("VELOCITY_SERVER is not set")
	}

	version := os.Getenv("VELOCITY_API_VERSION")
	if version == "" {
		version = "v1"
	}

	// test if it is in versions
	found := false
	for _, v := range versions {
		if string(v) == version {
			found = true
			break
		}
	}

	if !found {
		return nil, errors.New("VELOCITY_API_VERSION is not a valid version")
	}

	return &VelocityClient{
		BaseURL:    baseURL,
		APIVersion: Version(version),
	}, nil
}
