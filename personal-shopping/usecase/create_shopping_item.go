package usecase

type CreateShoppingItemUseCase interface {
	Execute(input CreateShoppingItemInput) (CreateShoppingItemOutput, error)
}

type createShoppingItemUseCase struct{}

func (c *createShoppingItemUseCase) Execute(input CreateShoppingItemInput) (CreateShoppingItemOutput, error) {
	return CreateShoppingItemOutput{
		ID:     1,
		Name:   input.Name,
		Status: "unpurchased",
		Category: struct {
			ID   int64
			Name string
		}{
			ID:   input.CategoryID,
			Name: "foods",
		},
	}, nil
}

type CreateShoppingItemInput struct {
	Name       string
	CategoryID int64
}

type CreateShoppingItemOutput struct {
	ID       int64
	Name     string
	Status   string
	Category struct {
		ID   int64
		Name string
	}
}

func NewCreateShoppingItemUseCase() CreateShoppingItemUseCase {
	return &createShoppingItemUseCase{}
}
