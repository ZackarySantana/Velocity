package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/zackarysantana/velocity/src/config"
	"github.com/zackarysantana/velocity/src/velocity"
)

// routeinestart2 is a version that is just the handler func

func (a *api[T]) routineStart(w http.ResponseWriter, r *http.Request) {
	body := velocity.APIStartRoutineRequest{}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err = body.Config.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if routine := body.Config.Routines.GetRoutine(body.Routine); routine == nil {
		http.Error(w, fmt.Sprintf("'%s' routine not found", body.Routine), http.StatusBadRequest)
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
	w.WriteHeader(http.StatusOK)
}
