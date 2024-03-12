package handlers

import "net/http"

func Readiness(w http.ResponseWriter, r *http.Request) {
	RespondWithCode(w, http.StatusOK)
}

func Liveness(w http.ResponseWriter, r *http.Request) {
	RespondWithCode(w, http.StatusOK)
}

func RespondWithCode(w http.ResponseWriter, code int) {
	w.WriteHeader(code)
}
