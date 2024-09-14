package velocity

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/samber/oops"
)

type baseClient struct {
	client *http.Client
	base   string

	headers map[string]string
}

func newBaseClient(base string) *baseClient {
	return &baseClient{
		client: &http.Client{},
		base:   base,
	}
}

func (c *baseClient) do(ctx context.Context, method, path string, payload interface{}) (*http.Response, error) {
	var body io.Reader
	if payload != nil {
		bodyBuffer := &bytes.Buffer{}
		err := json.NewEncoder(bodyBuffer).Encode(payload)
		if err != nil {
			return nil, err
		}
		body = bodyBuffer
	}

	req, err := http.NewRequestWithContext(ctx, method, c.base+path, body)
	if err != nil {
		return nil, err
	}
	for k, v := range c.headers {
		req.Header.Set(k, v)
	}
	if payload != nil && req.Header.Get("Content-Type") == "" {
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

func (c *baseClient) doAndDecode(ctx context.Context, method, path string, payload, dest interface{}) (*http.Response, error) {
	resp, err := c.do(ctx, method, path, payload)
	if err != nil {
		return resp, err
	}
	defer resp.Body.Close()
	return resp, json.NewDecoder(resp.Body).Decode(dest)
}
