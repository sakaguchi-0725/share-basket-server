package main

import (
	"fmt"
	"log"
	"sharebasket/core"
	"sharebasket/infra/db"
	"sharebasket/presentation/server"
	"sharebasket/registry"
)

func main() {
	cfg, err := core.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	conn, err := db.New(cfg.DB)
	if err != nil {
		log.Fatal(err)
	}

	repo, err := registry.NewRepository(conn, cfg)
	if err != nil {
		log.Fatalf("failed to initialize repository: %v", err)
	}

	service := registry.NewService(repo)
	usecase := registry.NewUseCase(repo, service)

	logger := core.NewLogger(cfg.Env)

	srv := server.New(fmt.Sprintf(":%v", cfg.Port))
	srv.MapHandler(usecase, logger)
	srv.Run()
}
