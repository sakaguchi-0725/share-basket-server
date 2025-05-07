package main

import (
	"log"
	"share-basket-server/core/config"
	"share-basket-server/presentation/server"
	"share-basket-server/registry"
)

func main() {
	cfg := config.Load()

	handlers, err := registry.Inject(cfg)
	if err != nil {
		log.Fatal(err)
	}

	s := server.New(":8080")
	s.MapRoutes(cfg.FrontendURL, handlers)

	s.Run()
}
