package response

import (
	"encoding/json"
	"net/http"
)

func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

func StatusOK[T any](w http.ResponseWriter, body T) error {
	return httpJSON(w, http.StatusOK, body)
}

func httpJSON[T any](w http.ResponseWriter, status int, body T) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(body); err != nil {
		return err
	}

	return nil
}
