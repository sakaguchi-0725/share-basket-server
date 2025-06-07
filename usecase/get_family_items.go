package usecase

import (
	"context"
	"sharebasket/domain/model"
	"sharebasket/domain/repository"
)

type (
	GetFamilyItems interface {
		Execute(ctx context.Context, in GetFamilyItemsInput) ([]GetFamilyItemOutput, error)
	}

	GetFamilyItemsInput struct {
		UserID   string
		Status   string
		FamilyID string
	}

	GetFamilyItemOutput struct {
		ID         int64
		Name       string
		Status     string
		CategoryID int64
		FamilyID   string
	}

	getFamilyItems struct {
		familyItemRepo repository.FamilyItem
	}
)

func (g *getFamilyItems) Execute(ctx context.Context, in GetFamilyItemsInput) ([]GetFamilyItemOutput, error) {
	familyID, err := model.ParseFamilyID(in.FamilyID)
	if err != nil {
		return nil, err
	}

	status, err := in.GetStatus()
	if err != nil {
		return nil, err
	}

	familyItems, err := g.familyItemRepo.Get(ctx, familyID, status)
	if err != nil {
		return nil, err
	}

	return g.makeOutput(familyItems), nil
}

func (g *getFamilyItems) makeOutput(items []model.FamilyItem) []GetFamilyItemOutput {
	outputs := make([]GetFamilyItemOutput, len(items))
	for i, item := range items {
		outputs[i] = GetFamilyItemOutput{
			ID:         *item.ID,
			Name:       item.Name,
			Status:     item.Status.String(),
			CategoryID: item.CategoryID,
			FamilyID:   item.FamilyID.String(),
		}
	}

	return outputs
}

func (g *GetFamilyItemsInput) GetStatus() (*model.ShoppingStatus, error) {
	var status model.ShoppingStatus
	if g.Status == "" {
		return nil, nil
	}

	status, err := model.ParseShoppingStatus(g.Status)
	if err != nil {
		return nil, err
	}

	return &status, nil
}

func NewGetFamilyItems(f repository.FamilyItem) GetFamilyItems {
	return &getFamilyItems{
		familyItemRepo: f,
	}
}
