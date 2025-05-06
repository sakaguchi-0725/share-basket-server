package domain

import (
	"errors"
)

type ShoppingStatus string

const (
	UnPurchased ShoppingStatus = "UnPurchased"
	InTheCart   ShoppingStatus = "InTheCart"
	Purchased   ShoppingStatus = "Purchased"
)

var (
	ErrInvalidShoppingStatus = errors.New("invalid shopping status")
)

func NewShoppingStatus(s string) (ShoppingStatus, error) {
	switch s {
	case UnPurchased.String():
		return UnPurchased, nil
	case InTheCart.String():
		return InTheCart, nil
	case Purchased.String():
		return Purchased, nil
	default:
		return "", ErrInvalidShoppingStatus
	}
}

func (status ShoppingStatus) String() string {
	return string(status)
}
