package model

import (
	"errors"
	"sharebasket/core"
)

type FamilyItem struct {
	ID         *int64
	Name       string
	Status     ShoppingStatus
	CategoryID int64
	FamilyID   FamilyID
	CreatedBy  AccountID
}

func NewFamilyItem(
	name string,
	status *ShoppingStatus,
	categoryID int64,
	familyID FamilyID,
	createdBy AccountID,
) (FamilyItem, error) {
	if name == "" {
		return FamilyItem{}, core.NewInvalidError(errors.New("family shopping item name is required"))
	}

	if categoryID <= 0 {
		return FamilyItem{}, core.NewInvalidError(errors.New("category ID is required"))
	}

	// ステータスが指定されていない場合、UnPurchasedにする
	shoppingStatus := UnPurchased
	if status != nil {
		shoppingStatus = core.Derefer(status)
	}

	return FamilyItem{
		Name:       name,
		Status:     shoppingStatus,
		CategoryID: categoryID,
		FamilyID:   familyID,
		CreatedBy:  createdBy,
	}, nil
}
