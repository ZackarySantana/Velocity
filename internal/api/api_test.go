package api_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type RequestOptions struct {
	method string
	path   string

	headers map[string]string

	body        io.Reader
	hasJSONBody bool
}

type RequestOptionFunc func(*RequestOptions)

func WithJSONBody(body interface{}) RequestOptionFunc {
	return func(r *RequestOptions) {
		r.headers["Content-Type"] = "application/json"
		bodyAsJsonString, err := json.Marshal(body)
		if err != nil {
			panic(err)
		}
		r.body = bytes.NewBuffer(bodyAsJsonString)
		r.hasJSONBody = true
	}
}

func IsAdminRequest(isAdmin bool) RequestOptionFunc {
	return func(r *RequestOptions) {
		if isAdmin {
			r.headers["Username"] = "admin"
			r.headers["Password"] = "admin"
		} else {
			r.headers["Username"] = "user"
			r.headers["Password"] = "user"
		}
	}
}

// Let's use configurable functions
func NewRequest(method, path string, opts ...RequestOptionFunc) *http.Request {
	defaultOpts := &RequestOptions{
		method:      method,
		path:        path,
		headers:     map[string]string{},
		body:        nil,
		hasJSONBody: false,
	}

	for _, opt := range opts {
		opt(defaultOpts)
	}

	req, err := http.NewRequest(defaultOpts.method, defaultOpts.path, defaultOpts.body)
	if err != nil {
		panic(err)
	}

	for k, v := range defaultOpts.headers {
		req.Header.Set(k, v)
	}

	return req
}
