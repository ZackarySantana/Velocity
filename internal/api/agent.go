package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/zackarysantana/velocity/src/velocity"
)

func (a *api) agentGetTask() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		test, err := a.repository.Test.Load(r.Context(), id).Unwrap()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Println(id, "test", test)
		if test == nil {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		resp := velocity.AgentGetTestResponse{Test: *test}
		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}
