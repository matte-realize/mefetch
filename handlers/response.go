package handlers

import (
	"encoding/json"
	"net/http"
)

func errorResponse(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(struct {
		Error string `json:"error"`
	} {
		Error: message,
	})
}

func successResponse(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

func handleError(w http.ResponseWriter, err error, message string, status int) bool {
	if err != nil {
		errorResponse(w, message, status)
		return true
	}
	
	return false
}