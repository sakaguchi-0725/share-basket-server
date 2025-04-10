package registry

import (
	"share-basket/personal-shopping/presentation/handler"
	"share-basket/personal-shopping/usecase"
)

type container struct {
	createShoppingItem usecase.CreateShoppingItemUseCase
	getShoppingItems   usecase.GetShoppingItemsUseCase
	getAccount         usecase.GetAccountUseCase
}

func Inject() (*container, error) {
	getShoppingItems := usecase.NewGetShoppingItemsUseCase()
	getAccount := usecase.NewGetAccountUseCase()
	createShoppingItem := usecase.NewCreateShoppingItemUseCase()

	return &container{
		createShoppingItem,
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

func (c *container) CreateShoppingItemHandler() handler.CreateShoppingItemHandler {
	return handler.NewCreateShoppingItemHandler(c.createShoppingItem)
}
