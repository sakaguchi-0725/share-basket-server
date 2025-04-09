package main

import (
	"log"
	"share-basket/personal-shopping/presentation/server"
	"share-basket/personal-shopping/registry"
)

func main() {
	container, err := registry.Inject()
	if err != nil {
		log.Fatalf("failed to inject dependencies: %v", err)
	}

	server := server.New(":50051")
	services := registry.NewServices(container)
	server.MapServices(services)

	server.Run()
}
