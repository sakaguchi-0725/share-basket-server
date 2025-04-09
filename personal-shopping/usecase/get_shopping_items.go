package usecase

type GetShoppingItemsUseCase interface {
	Execute(status string) ([]GetShoppingItemOutput, error)
}

type getShoppingItemsUseCase struct{}

func (usecase *getShoppingItemsUseCase) Execute(status string) ([]GetShoppingItemOutput, error) {
	// TODO: implements.
	return []GetShoppingItemOutput{
		{
			ID:     1,
			Name:   "牛乳",
			Status: "unpurchased",
			Category: GetShoppingItemCategory{
				ID:   1,
				Name: "foods",
			},
		},
		{
			ID:     2,
			Name:   "卵",
			Status: "unpurchased",
			Category: GetShoppingItemCategory{
				ID:   1,
				Name: "foods",
			},
		},
	}, nil
}

type GetShoppingItemOutput struct {
	ID       int64
	Name     string
	Status   string
	Category GetShoppingItemCategory
}

type GetShoppingItemCategory struct {
	ID   int64
	Name string
}

func NewGetShoppingItemsUseCase() GetShoppingItemsUseCase {
	return &getShoppingItemsUseCase{}
}
