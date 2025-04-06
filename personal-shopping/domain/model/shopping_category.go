package model

type ShoppingCategory struct {
	ID   *int64
	Name string
}

func NewShoppingCategory(id *int64, name string) ShoppingCategory {
	return ShoppingCategory{
		ID:   id,
		Name: name,
	}
}
