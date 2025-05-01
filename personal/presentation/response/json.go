package response

import (
	"encoding/json"
	"net/http"
)

func NoContent(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func JSON[T any](w http.ResponseWriter, status int, body T) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(body); err != nil {
		return err
	}

	return nil
}
