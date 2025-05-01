package model

type ShoppingCategory struct {
	ID   uint
	Name string
}

func NewShoppingCategory(id uint, name string) ShoppingCategory {
	return ShoppingCategory{
		ID:   id,
		Name: name,
	}
}
