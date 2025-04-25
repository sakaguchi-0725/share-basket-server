package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Handlers struct {
	PingHandler http.HandlerFunc
}

func RegisterRoutes(r chi.Router, handlers Handlers) {
	r.Route("/personal", func(r chi.Router) {
		r.Get("/ping", handlers.PingHandler)
	})
}
