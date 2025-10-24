package api

import (
	"encoding/json"
	"net/http"

	"github.com/zackarysantana/velocity/src/config"
	"github.com/zackarysantana/velocity/src/velocity"
)

func (a *api[T]) routineStart(w http.ResponseWriter, r *http.Request) {
	body := velocity.APIStartRoutineRequest{}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err = body.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ec, err := config.CreateEntity(&body.Config, config.CreateEntityOptions[T]{
		Id:              a.idCreator,
		FilterToRoutine: body.Routine,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err = a.service.StartRoutine(r.Context(), ec, body.Routine); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resp := velocity.APIStartRoutineResponse{Id: ec.Routines[0].Id}
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (a *api[T]) routineGet(w http.ResponseWriter, r *http.Request) {
	id, err := a.idCreator.Read(r.PathValue("id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	routines, err := a.repository.Routine().Load(r.Context(), []T{id})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(routines) == 0 {
		http.Error(w, "routine not found", http.StatusNotFound)
		return
	}
	resp := velocity.APIGetRoutineResponse[T]{Routine: *routines[0]}
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
