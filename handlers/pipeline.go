package handlers

import "net/http"

type step func() error

func run(w http.ResponseWriter, steps ...step) bool {
	for _, s := range steps {
		if err := s(); err != nil {
			errorResponse(w, err.Error(), http.StatusInternalServerError)
			return false
		}
	}

	return true
}