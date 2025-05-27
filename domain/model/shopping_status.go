package model

import (
	"errors"
	"sharebasket/core"
)

type ShoppingStatus string

const (
	UnPurchased ShoppingStatus = "UnPurchased"
	InTheCart   ShoppingStatus = "InTheCart"
	Purchased   ShoppingStatus = "Purchased"
)

func ParseShoppingStatus(s string) (ShoppingStatus, error) {
	if s == UnPurchased.String() || s == InTheCart.String() || s == Purchased.String() {
		return ShoppingStatus(s), nil
	}

	return "", core.NewInvalidError(errors.New("invalid shopping status"))
}

func (s ShoppingStatus) String() string {
	return string(s)
}
