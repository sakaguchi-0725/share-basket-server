package handler

import (
	"net/http"
	"share-basket-server/personal/presentation/response"
)

func MakePingHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response.StatusOK(w, map[string]string{
			"message": "pong",
		})
	}
}
