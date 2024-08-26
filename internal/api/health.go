package api

import "net/http"

func (a *api[T]) health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
