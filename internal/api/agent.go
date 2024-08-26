package api

import (
	"encoding/json"
	"net/http"

	"github.com/zackarysantana/velocity/src/velocity"
)

func (a *api[T]) agentGetTask(w http.ResponseWriter, r *http.Request) {
	pathId := r.PathValue("id")
	id, err := a.idCreator.Read(pathId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	tests, err := a.repository.Test.Load(r.Context(), []T{id})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if tests == nil || len(tests) == 0 {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	resp := velocity.AgentGetTestResponse[T]{Test: *tests[0]}
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
