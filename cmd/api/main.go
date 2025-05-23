package main

import (
	"log"
	"sharebasket/core"
	"sharebasket/infra/db"
	"sharebasket/presentation/server"
	"sharebasket/registry"
)

func main() {
	conn, err := db.New()
	if err != nil {
		log.Fatal(err)
	}

	repo, err := registry.NewRepository(conn)
	if err != nil {
		log.Fatalf("failed to initialize repository: %v", err)
	}

	service := registry.NewService(repo)
	usecase := registry.NewUseCase(repo, service)

	logger := core.NewLogger()

	srv := server.New(8080)
	srv.MapHandler(usecase, logger)
	srv.Run()
}
