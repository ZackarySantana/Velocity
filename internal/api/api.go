package api

import (
	"fmt"
	"net/http"

	"github.com/zackarysantana/velocity/internal/service"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("api")

var (
	middlewares    = []func(http.Handler) http.Handler{}
	apiMiddlewares = []func(http.Handler) http.Handler{
		func(h http.Handler) http.Handler {
			return otelhttp.NewHandler(h, "api")
		},
		func(h http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Println("api middleware")
				h.ServeHTTP(w, r)
			})
		},
	}
	agentMiddlewares = []func(http.Handler) http.Handler{
		func(h http.Handler) http.Handler {
			return otelhttp.NewHandler(h, "agent-api")
		},
		func(h http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Println("agent middleware")
				h.ServeHTTP(w, r)
			})
		},
	}
)

type api struct {
	service   service.Service
	idCreator service.IdCreator
}

func New(service service.Service, idCreator service.IdCreator) http.Handler {
	a := &api{service: service, idCreator: idCreator}

	rootMux := http.NewServeMux()

	apiMux := http.NewServeMux()
	apiMux.Handle("GET /health", a.health())
	apiMux.Handle("POST /routine/start", a.routineStart())
	rootMux.Handle("/", applyMiddleware(apiMux, apiMiddlewares...))

	agentMux := http.NewServeMux()
	agentMux.Handle("GET /health", a.health())
	rootMux.Handle("/agent/", http.StripPrefix("/agent", applyMiddleware(agentMux, agentMiddlewares...)))

	return applyMiddleware(rootMux, middlewares...)
}

func applyMiddleware(h http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for _, m := range middlewares {
		h = m(h)
	}
	return h
}
