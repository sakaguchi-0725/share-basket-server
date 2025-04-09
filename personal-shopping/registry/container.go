package registry

import (
	"share-basket/personal-shopping/presentation/handler"
	"share-basket/personal-shopping/usecase"
)

type container struct {
	getShoppingItems usecase.GetShoppingItemsUseCase
}

func Inject() (*container, error) {
	getShoppingItems := usecase.NewGetShoppingItemsUseCase()

	return &container{
		getShoppingItems,
	}, nil
}

func (c *container) GetShoppingItemsHandler() handler.GetShoppingItemsHandler {
	return handler.NewGetShoppingItemsHandler(c.getShoppingItems)
}
