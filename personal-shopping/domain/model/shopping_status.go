package model

import "errors"

type ShoppingStatus string

const (
	UnPurchased ShoppingStatus = "un_purchased"
	Purchased   ShoppingStatus = "purchased"
)

func NewShoppingStatus(s string) (ShoppingStatus, error) {
	switch s {
	case UnPurchased.String():
		return UnPurchased, nil
	case Purchased.String():
		return Purchased, nil
	default:
		return "", errors.New("invalid shopping status")
	}
}

func (status ShoppingStatus) String() string {
	return string(status)
}
