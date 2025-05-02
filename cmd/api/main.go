package main

import (
	"log"
	"log/slog"
	"share-basket-server/core/config"
	"share-basket-server/core/db"
	"share-basket-server/core/logger"
	"share-basket-server/core/server"
	personalRouter "share-basket-server/personal/presentation/router"
	personalRegistry "share-basket-server/personal/registry"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func main() {
	cfg := config.Load()
	logger := logger.New(cfg.Env)
	slog.SetDefault(logger.Logger)

	db, err := db.New(cfg.DB)
	if err != nil {
		log.Fatal(err)
	}

	personalHandlers, err := personalRegistry.Inject(db, cfg.AWS)
	if err != nil {
		log.Fatal(err)
	}

	s := server.New(":8080")
	s.MapRoutes(func(r chi.Router) {
		r.Use(middleware.Logger)
		r.Use(middleware.Recoverer)
		r.Use(cors.Handler(cors.Options{
			AllowedOrigins:   []string{cfg.FrontendURL},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			AllowCredentials: true,
		}))

		personalRouter.RegisterRoutes(r, personalHandlers)
	})

	s.Run()
}
