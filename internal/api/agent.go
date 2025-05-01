package api

import (
	"encoding/json"
	"net/http"

	"github.com/zackarysantana/velocity/internal/service"
	"github.com/zackarysantana/velocity/src/velocity"
)

func (a *api[T]) agentGetTask(w http.ResponseWriter, r *http.Request) {
	pathId := r.PathValue("id")
	id, err := a.idCreator.Read(pathId)
	if err != nil {
		http.Error(w, "incorrectly formatted id", http.StatusNotFound)
		return
	}
	tests, err := a.repository.Test().Load(r.Context(), []T{id})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(tests) == 0 {
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

// h(agentMux, "POST /priority_queue/pop", a.agentGetTask)
// h(agentMux, "POST /priority_queue/done", a.agentGetTask)
// h(agentMux, "POST /priority_queue/unfinished", a.agentGetTask)

func (a *api[T]) agentPriorityQueuePop(w http.ResponseWriter, r *http.Request) {
	var req velocity.AgentPopRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	item, err := a.testQueue.Pop(r.Context(), req.Type)
	if err != nil {
		if err == service.ErrEmptyQueue {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resp := velocity.AgentPopResponse{
		Popped: service.PriorityQueuePoppedItem[any, any]{
			Id:        item.Id,
			Payload:   item.Payload,
			CreatedOn: item.CreatedOn,
		},
	}
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
