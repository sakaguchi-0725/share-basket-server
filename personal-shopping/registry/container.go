package registry

import (
	"share-basket/personal-shopping/presentation/handler"
	"share-basket/personal-shopping/usecase"
)

type container struct {
	getShoppingItems usecase.GetShoppingItemsUseCase
	getAccount       usecase.GetAccountUseCase
}

func Inject() (*container, error) {
	getShoppingItems := usecase.NewGetShoppingItemsUseCase()
	getAccount := usecase.NewGetAccountUseCase()

	return &container{
		getShoppingItems,
		getAccount,
	}, nil
}

func (c *container) GetShoppingItemsHandler() handler.GetShoppingItemsHandler {
	return handler.NewGetShoppingItemsHandler(c.getShoppingItems)
}

func (c *container) GetAccountHandler() handler.GetAccountHandler {
	return handler.NewGetAccountHandler(c.getAccount)
}
