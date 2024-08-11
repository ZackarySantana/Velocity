package api

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/zackarysantana/velocity/internal/service"
	"github.com/zackarysantana/velocity/src/writer"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("api")

var ()

type api struct {
	repository *service.Repository
	service    service.Service
	idCreator  service.IdCreator
	logger     *slog.Logger
}

func New(repository *service.Repository, service service.Service, idCreator service.IdCreator, logger *slog.Logger) http.Handler {
	a := &api{service: service, idCreator: idCreator, repository: repository, logger: logger}

	middlewares := []func(http.Handler) http.Handler{
		func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Start time
				start := time.Now()

				// Capture the response status code
				rw := &writer.Response{ResponseWriter: w}

				// Call the next handler
				next.ServeHTTP(rw, r)

				method := r.Method
				if method == "" {
					method = "GET"
				}

				// Log the details
				a.logger.Info(
					r.URL.Path,
					"address",
					r.RemoteAddr,
					"status",
					fmt.Sprintf("%d", rw.StatusCode()),
					"method",
					method,
					"duration",
					time.Since(start).String(),
				)
			})
		},
	}
	apiMiddlewares := []func(http.Handler) http.Handler{
		func(next http.Handler) http.Handler {
			return otelhttp.NewHandler(next, "api")
		},
		func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Println("api middleware")
				next.ServeHTTP(w, r)
			})
		},
	}
	agentMiddlewares := []func(http.Handler) http.Handler{
		func(next http.Handler) http.Handler {
			return otelhttp.NewHandler(next, "agent-api")
		},
		func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Println("agent middleware")
				next.ServeHTTP(w, r)
			})
		},
	}

	rootMux := http.NewServeMux()

	apiMux := http.NewServeMux()
	apiMux.Handle("GET /health", a.health())
	apiMux.Handle("POST /routine/start", a.routineStart())
	rootMux.Handle("/", applyMiddleware(apiMux, apiMiddlewares...))

	agentMux := http.NewServeMux()
	agentMux.Handle("GET /health", a.health())
	agentMux.Handle("GET /test/{id}", a.agentGetTask())
	rootMux.Handle("/agent/", http.StripPrefix("/agent", applyMiddleware(agentMux, agentMiddlewares...)))

	return applyMiddleware(rootMux, middlewares...)
}

func applyMiddleware(h http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for _, m := range middlewares {
		h = m(h)
	}
	return h
}
