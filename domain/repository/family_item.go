package repository

import "sharebasket/domain/model"

type FamilyItem interface {
	Get(id model.FamilyID) ([]model.FamilyItem, error)
	Store(item *model.FamilyItem) error
}
