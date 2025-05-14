package handler

import (
	"net/http"
	"sharebasket/presentation/response"
)

func NewHealthCheck() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response.StatusOK(w, map[string]string{
			"message": "ok",
		})
	}
}
