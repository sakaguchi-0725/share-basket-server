package main

import (
	personalRouter "share-basket-server/personal/presentation/router"
	personalRegistry "share-basket-server/personal/registry"
	"share-basket-server/server"

	"github.com/go-chi/chi/v5"
)

func main() {
	s := server.New(":8080")

	personalHandlers := personalRegistry.Inject()
	s.MapRoutes(func(r chi.Router) {
		personalRouter.RegisterRoutes(r, personalHandlers)
	})

	s.Run()
}
