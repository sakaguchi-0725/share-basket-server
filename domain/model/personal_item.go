package model

import (
	"errors"
	"sharebasket/core"
)

type PersonalItem struct {
	ID         *int64
	Name       string
	Status     ShoppingStatus
	CategoryID int64
	AccountID  AccountID
}

func NewPersonalItem(name string, status *ShoppingStatus, categoryID int64, accID AccountID) (PersonalItem, error) {
	if name == "" {
		return PersonalItem{}, core.NewInvalidError(errors.New("personal shopping item name is required"))
	}

	if categoryID <= 0 {
		return PersonalItem{}, core.NewInvalidError(errors.New("category ID is required"))
	}

	// ステータスが指定されていない場合、UnPurchasedにする
	shoppingStatus := UnPurchased
	if status != nil {
		shoppingStatus = core.Derefer(status)
	}

	return PersonalItem{
		Name:       name,
		Status:     shoppingStatus,
		CategoryID: categoryID,
		AccountID:  accID,
	}, nil
}

// 指定されたアカウントがこのアイテムの所有者かを確認
func (p *PersonalItem) CheckOwner(accID AccountID) error {
	if p.AccountID != accID {
		return errors.New("you don't have permission to this item")
	}
	return nil
}

// Update は指定された値でアイテムを更新
// 空の値は無視され、既存の値が維持
func (p *PersonalItem) Update(name string, status *ShoppingStatus, categoryID *int64) error {
	if name != "" {
		p.Name = name
	}

	if status != nil {
		p.Status = *status
	}

	if categoryID != nil {
		if *categoryID <= 0 {
			return errors.New("category ID must be positive")
		}
		p.CategoryID = *categoryID
	}

	return nil
}
