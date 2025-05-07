//go:generate mockgen -destination=../test/mock/usecase/get_personal_shopping_items_input.go . GetPersonalShoppingItemsInputPort
//go:generate mockgen -destination=../test/mock/usecase/get_personal_shopping_items_output.go . GetPersonalShoppingItemsOutputPort
package usecase

import (
	"context"
	"errors"
	"share-basket-server/core/apperr"
	"share-basket-server/core/logger"
	"share-basket-server/core/util"
	"share-basket-server/domain"
)

type (
	GetPersonalShoppingItemsInputPort interface {
		Execute(ctx context.Context, input GetPersonalShoppingItemsInput, output GetPersonalShoppingItemsOutputPort) error
	}

	GetPersonalShoppingItemsInput struct {
		UserID string
		Status string
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
		logger       logger.Logger
	}
)

func (g *getPersonalShoppingItemsInteractor) Execute(ctx context.Context, input GetPersonalShoppingItemsInput, output GetPersonalShoppingItemsOutputPort) error {
	userID, err := domain.ParseUserID(input.UserID)
	if err != nil {
		g.logger.
			With("user id", input.UserID).
			With("error", err).
			Info("invalid user id")
		return apperr.NewInvalidError(err)
	}

	account, err := g.accountRepo.FindByUserID(userID)
	if err != nil {
		if errors.Is(err, apperr.ErrDataNotFound) {
			g.logger.
				With("user id", userID.String()).
				With("error", err).
				Info("account not found")
			return apperr.NewInvalidError(err)
		}

		g.logger.With("error", err).Error("failed to find account")
		return err
	}

	var status *domain.ShoppingStatus
	if input.Status != "" {
		s, err := domain.NewShoppingStatus(input.Status)
		if err != nil {
			g.logger.
				With("status", input.Status).
				With("error", err).
				Info("invalid status")
			return apperr.NewInvalidError(err)
		}
		status = util.Ptr(s)
	}

	items, err := g.personalRepo.GetAll(account.ID, status)
	if err != nil {
		g.logger.
			With("err", err).
			Error("failed to get personal shopping items")
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
	logger logger.Logger,
) GetPersonalShoppingItemsInputPort {
	return &getPersonalShoppingItemsInteractor{accountRepo, personalRepo, logger}
}
