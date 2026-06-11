package handlers

import (
	"encoding/json"
	"net/http"
)

func Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	resp := struct {
		Status   string `json:"status"`
		Message  string `json:"message"`
	} {
		Status:  "ok",
		Message: "mefetch is running",
	}

	json.NewEncoder(w).Encode(resp)
}