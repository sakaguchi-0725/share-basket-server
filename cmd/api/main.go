package main

import (
	"log"
	"log/slog"
	"share-basket-server/core/config"
	"share-basket-server/core/db"
	"share-basket-server/core/logger"
	"share-basket-server/presentation/server"
	"share-basket-server/registry"
)

func main() {
	cfg := config.Load()
	logger := logger.New(cfg.Env)
	slog.SetDefault(logger.Logger)

	db, err := db.New(cfg.DB)
	if err != nil {
		log.Fatal(err)
	}

	handlers, err := registry.Inject(db, cfg.AWS)
	if err != nil {
		log.Fatal(err)
	}

	s := server.New(":8080")
	s.MapRoutes(cfg.FrontendURL, handlers)

	s.Run()
}
