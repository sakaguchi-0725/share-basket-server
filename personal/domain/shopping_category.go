package domain

type (
	ShoppingCategory struct {
		ID   uint
		Name string
	}

	ShoppingCategoryRepository interface {
		GetAll() ([]ShoppingCategory, error)
	}
)

func NewShoppingCategory(id uint, name string) ShoppingCategory {
	return ShoppingCategory{
		ID:   id,
		Name: name,
	}
}
