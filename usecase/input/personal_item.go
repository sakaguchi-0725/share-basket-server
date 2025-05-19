package input

import (
	"sharebasket/core"
	"sharebasket/domain/model"
)

type (
	GetPersonalItem struct {
		UserID string
		Status string
	}

	CreatePersonalItem struct {
		UserID     string
		Name       string
		Status     string
		CategoryID int64
	}
)

func (g GetPersonalItem) ParseShoppingStatus() (*model.ShoppingStatus, error) {
	if g.Status == "" {
		return nil, nil
	}

	status, err := model.ParseShoppingStatus(g.Status)
	if err != nil {
		return nil, core.NewInvalidError(err)
	}

	return &status, nil
}
