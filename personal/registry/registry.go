package registry

import (
	"share-basket-server/personal/presentation/handler"
	"share-basket-server/personal/presentation/router"
)

func Inject() router.Handlers {
	return router.Handlers{
		PingHandler: handler.MakePingHandler(),
	}
}
