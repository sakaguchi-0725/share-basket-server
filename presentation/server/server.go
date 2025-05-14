package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sharebasket/core"
	"sharebasket/presentation/handler"
	customMiddleware "sharebasket/presentation/middleware"
	"sharebasket/registry"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	addr   string
	router *chi.Mux
	server *http.Server
}

func New(addr uint) *Server {
	r := chi.NewRouter()
	addrStr := fmt.Sprintf(":%v", addr)

	return &Server{
		addr:   addrStr,
		router: r,
		server: &http.Server{
			Addr:         addrStr,
			Handler:      r,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  5 * time.Second,
		},
	}
}

func (s *Server) MapHandler(usecase registry.UseCase, logger core.Logger) {
	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Recoverer)

	s.router.Get("/health-check", handler.NewHealthCheck())
	s.router.Post("/login", handler.NewLogin(usecase.NewLogin(), logger))
	s.router.Post("/signup", handler.NewSignUp(usecase.NewSignUp(), logger))
	s.router.Post("/signup/confirm", handler.NewSignUpConfirm(usecase.NewSignUpConfirm(), logger))
	s.router.Post("/logout", handler.NewLogout())

	s.router.Group(func(r chi.Router) {
		r.Use(customMiddleware.Auth(usecase.NewVerifyToken(), logger))
		r.Get("/me", handler.NewGetAccount(usecase.NewGetAccount(), logger))
		r.Get("/categories", handler.NewGetCategories(usecase.NewGetCategories()))

		r.Route("/personal", func(r chi.Router) {
			r.Get("/items", handler.NewGetPersonalItems(usecase.NewGetPersonalItems(), logger))
			r.Post("/items", handler.NewCreatePersonalItem(usecase.NewCreatePersonalItem(), logger))
			r.Put("/items/{id}", handler.NewUpdatePersonalItem(usecase.NewUpdatePersonalItem(), logger))
			r.Delete("/items/{id}", handler.NewDeletePersonalItem(usecase.NewDeletePersonalItem(), logger))
		})
	})
}

func (s *Server) Run() {
	go func() {
		log.Println("HTTP server is running on", s.addr)
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to serve: %v\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		log.Fatalf("server shutdown failed: %v", err)
	}
	log.Println("HTTP server shutdown completed")
}
