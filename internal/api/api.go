package api

import (
	"net/http"

	"github.com/zackarysantana/velocity/internal/service"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("api")

var middlewares = []func(http.Handler) http.Handler{}

type api struct {
	service service.Service
}

func New(service service.Service) *http.ServeMux {
	mux := http.NewServeMux()
	a := &api{service: service}

	handlers := map[string]http.Handler{}
	handlers["GET /health"] = a.health()

	for r, h := range handlers {
		handler := http.Handler(h)
		for _, m := range middlewares {
			handler = m(handler)
		}
		handler = otelhttp.NewHandler(handler, r)

		mux.Handle(r, handler)
	}

	return mux
}
