package router

import (
	"net/http"
	"share-basket-server/personal/presentation/middleware"

	"github.com/go-chi/chi/v5"
)

type Handlers struct {
	PingHandler          http.HandlerFunc
	SignUpHandler        http.HandlerFunc
	SignUpConfirmHandler http.HandlerFunc
	LoginHandler         http.HandlerFunc
}

func RegisterRoutes(r chi.Router, handlers Handlers) {
	r.Route("/personal", func(r chi.Router) {
		r.Use(middleware.WithError)
		r.Get("/ping", handlers.PingHandler)
		r.Post("/signup", handlers.SignUpHandler)
		r.Post("/signup/confirm", handlers.SignUpConfirmHandler)
		r.Post("/login", handlers.LoginHandler)
	})
}
