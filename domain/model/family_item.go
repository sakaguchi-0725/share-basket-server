package model

type FamilyItem struct {
	ID         *int64
	Name       string
	Status     ShoppingStatus
	CategoryID int64
	CreatedBy  AccountID
}
