package main

import (
	"log/slog"
	"net/http"

	"github.com/zackarysantana/velocity/internal/api"
)

func main() {
	mux := api.New(nil)

	slog.Info("Starting server", "addr", "0.0.0.0:8080")
	http.ListenAndServe("0.0.0.0:8080", mux)
}
