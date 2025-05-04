package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type (
	Server struct {
		addr   string
		router *chi.Mux
		server *http.Server
	}

	Handlers struct {
		PingHandler                  http.HandlerFunc
		SignUpHandler                http.HandlerFunc
		SignUpConfirmHandler         http.HandlerFunc
		LoginHandler                 http.HandlerFunc
		GetShoppingCaterogiesHandler http.HandlerFunc
	}
)

func New(addr string) *Server {
	r := chi.NewRouter()

	return &Server{
		addr:   addr,
		router: r,
		server: &http.Server{
			Addr:         addr,
			Handler:      r,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  5 * time.Second,
		},
	}
}

func (s *Server) MapRoutes(frontendURL string, handlers Handlers) {
	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Recoverer)
	s.router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{frontendURL},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
	}))

	s.router.Get("/ping", handlers.PingHandler)
	s.router.Post("/signup", handlers.SignUpHandler)
	s.router.Post("/signup/confirm", handlers.SignUpConfirmHandler)
	s.router.Post("/login", handlers.LoginHandler)
	s.router.Get("/categories", handlers.GetShoppingCaterogiesHandler)
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
