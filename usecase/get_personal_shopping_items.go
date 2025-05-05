//go:generate mockgen -destination=../mock/usecase/get_personal_shopping_items_input.go . GetPersonalShoppingItemsInputPort
//go:generate mockgen -destination=../mock/usecase/get_personal_shopping_items_output.go . GetPersonalShoppingItemsOutputPort
package usecase

import (
	"context"
	"errors"
	"share-basket-server/core/apperr"
	contextKey "share-basket-server/core/context"
	"share-basket-server/core/util"
	"share-basket-server/domain"
)

type (
	GetPersonalShoppingItemsInputPort interface {
		Execute(ctx context.Context, status string, output GetPersonalShoppingItemsOutputPort) error
	}

	GetPersonalShoppingItemsOutputPort interface {
		Render(ctx context.Context, output []GetPersonalShoppingItemOutput) error
	}

	GetPersonalShoppingItemOutput struct {
		ID         uint
		Name       string
		CategoryID uint
		Status     string
	}

	getPersonalShoppingItemsInteractor struct {
		accountRepo  domain.AccountRepository
		personalRepo domain.PersonalShoppingItemRepository
	}
)

func (g *getPersonalShoppingItemsInteractor) Execute(ctx context.Context, status string, output GetPersonalShoppingItemsOutputPort) error {
	userID := ctx.Value(contextKey.UserID).(domain.UserID)

	account, err := g.accountRepo.FindByUserID(userID)
	if err != nil {
		if errors.Is(err, apperr.ErrDataNotFound) {
			return apperr.NewInvalidError(err)
		}
		return err
	}

	var shoppingStatus *domain.ShoppingStatus
	if status != "" {
		s, err := domain.NewShoppingStatus(status)
		if err != nil {
			return apperr.NewInvalidError(err)
		}
		shoppingStatus = util.Ptr(s)
	}

	items, err := g.personalRepo.GetAll(account.ID, shoppingStatus)
	if err != nil {
		return err
	}

	return output.Render(ctx, g.makeOutputs(items))
}

func (g *getPersonalShoppingItemsInteractor) makeOutputs(items []domain.PersonalShoppingItem) []GetPersonalShoppingItemOutput {
	outputs := make([]GetPersonalShoppingItemOutput, len(items))

	for i, v := range items {
		outputs[i] = GetPersonalShoppingItemOutput{
			ID:         util.Derefer(v.ID),
			Name:       v.Name,
			CategoryID: v.CategoryID,
			Status:     v.Status.String(),
		}
	}

	return outputs
}

func NewGetPersonalShoppingItemsInteractor(
	accountRepo domain.AccountRepository,
	personalRepo domain.PersonalShoppingItemRepository,
) GetPersonalShoppingItemsInputPort {
	return &getPersonalShoppingItemsInteractor{accountRepo, personalRepo}
}
