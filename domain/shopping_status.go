package domain

type ShoppingStatus string

const (
	UnPurchased ShoppingStatus = "UnPurchased"
	InTheCart   ShoppingStatus = "InTheCart"
	Purchased   ShoppingStatus = "Purchased"
)

func (status ShoppingStatus) String() string {
	return string(status)
}
