package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/zackarysantana/velocity/src/velocity"
)

func (a *api) routineStart() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body := velocity.StartRoutineRequst{}
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
		ec := body.Config.ToEntity(a.idCreator)
		if err = a.service.StartRoutine(r.Context(), ec, body.Routine); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		resp := velocity.StartRoutineResponse{Id: ec.Routines[0].Id}
		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}
