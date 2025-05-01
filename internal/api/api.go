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

type api[ID any] struct {
	idCreator  service.IDCreator[ID]
	repository service.RepositoryManager[ID]
	service    service.Service[ID]

	testQueue service.PriorityQueue[ID, ID]

	logger *slog.Logger
}

// New creates a new http.Handler that serves the API. The given type is the
// type of the ids for data.
func New[T any](idCreator service.IDCreator[T], repository service.RepositoryManager[T], service service.Service[T], testQueue service.PriorityQueue[T, T], logger *slog.Logger) http.Handler {
	a := &api[T]{idCreator: idCreator, repository: repository, service: service, testQueue: testQueue, logger: logger}

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

	h := func(mux *http.ServeMux, pattern string, handler func(http.ResponseWriter, *http.Request)) {
		mux.Handle(pattern, otelhttp.NewHandler(http.HandlerFunc(handler), pattern))
	}

	rootMux := http.NewServeMux()

	apiMux := http.NewServeMux()
	h(apiMux, "GET /health", a.health)
	h(apiMux, "POST /routine/start", a.routineStart)
	rootMux.Handle("/", applyMiddleware(apiMux, apiMiddlewares...))

	agentMux := http.NewServeMux()
	h(agentMux, "GET /health", a.health)
	h(agentMux, "GET /test/{id}", a.agentGetTask)
	h(agentMux, "POST /priority_queue/pop", a.agentPriorityQueuePop)
	h(agentMux, "POST /priority_queue/done", a.agentGetTask)
	h(agentMux, "POST /priority_queue/unfinished", a.agentGetTask)
	rootMux.Handle("/agent/", http.StripPrefix("/agent", applyMiddleware(agentMux, agentMiddlewares...)))

	return applyMiddleware(rootMux, middlewares...)
}

func applyMiddleware(h http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for _, m := range middlewares {
		h = m(h)
	}
	return h
}
