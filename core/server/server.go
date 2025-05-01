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
)

type Server struct {
	addr   string
	router *chi.Mux
	server *http.Server
}

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

func (s *Server) MapRoutes(fn func(r chi.Router)) {
	fn(s.router)
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
