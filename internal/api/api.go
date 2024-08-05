package api

import (
	"net/http"

	"github.com/zackarysantana/velocity/internal/service"
	"github.com/zackarysantana/velocity/src/config/id"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("api")

var middlewares = []func(http.Handler) http.Handler{}

type api struct {
	service   service.Service
	idCreator id.Creator
}

func New(service service.Service, idCreator id.Creator) *http.ServeMux {
	mux := http.NewServeMux()
	a := &api{service: service, idCreator: idCreator}

	handlers := map[string]http.Handler{}
	handlers["GET /health"] = a.health()
	handlers["POST /routine/start"] = a.routineStart()

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
