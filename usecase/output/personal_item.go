package output

import (
	"sharebasket/core"
	"sharebasket/domain/model"
)

type GetPersonalItem struct {
	ID         int64
	Name       string
	Status     string
	CategoryID int64
}

func NewGetPersonalItemOutputs(items []model.PersonalItem) []GetPersonalItem {
	outputs := make([]GetPersonalItem, len(items))

	for i, v := range items {
		outputs[i] = GetPersonalItem{
			ID:         core.Derefer(v.ID),
			Name:       v.Name,
			Status:     v.Status.String(),
			CategoryID: v.CategoryID,
		}
	}

	return outputs
}
