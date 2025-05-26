package dao

import (
	"context"
	"sharebasket/core"
	"sharebasket/domain/model"
	"sharebasket/domain/repository"
	"sharebasket/infra/db"
	"time"
)

type (
	familyItemDao struct {
		conn *db.Conn
	}

	familyItemDto struct {
		ID         int64       `gorm:"primaryKey autoIncrement"`
		Name       string      `gorm:"not null"`
		Status     string      `gorm:"not null"`
		CategoryID int64       `gorm:"not null"`
		Category   categoryDto `gorm:"foreignKey:CategoryID;references:ID"`
		CreatedBy  string      `gorm:"not null"`
		Account    accountDto  `gorm:"foreignKey:CreatedBy;references:ID"`
		FamilyID   string      `gorm:"not null"`
		Family     familyDto   `gorm:"foreignKey:FamilyID;references:ID"`
		CreatedAt  time.Time   `gorm:"autoCreateTime"`
		UpdatedAt  time.Time   `gorm:"autoUpdateTime"`
	}

	familyItemDtos []familyItemDto
)

func (f *familyItemDao) Get(ctx context.Context, id model.FamilyID, status *model.ShoppingStatus) ([]model.FamilyItem, error) {
	var dtos familyItemDtos

	query := f.conn.Where("family_id = ?", id.String())

	if status != nil {
		query = query.Where("status = ?", status.String())
	}

	if err := query.Find(&dtos).Error; err != nil {
		return []model.FamilyItem{}, err
	}

	return dtos.toModels(), nil
}

func (f *familyItemDao) Store(ctx context.Context, item *model.FamilyItem) error {
	dto := newFamilyItemDto(item)

	if err := f.conn.Save(&dto).Error; err != nil {
		return err
	}

	if item.ID == nil {
		item.ID = core.Ptr(dto.ID)
	}

	return nil
}

func newFamilyItemDto(item *model.FamilyItem) familyItemDto {
	return familyItemDto{
		ID:     core.Derefer(item.ID),
		Name:   item.Name,
		Status: item.Status.String(),
	}
}

func (dtos familyItemDtos) toModels() []model.FamilyItem {
	items := make([]model.FamilyItem, len(dtos))

	for i, v := range dtos {
		items[i] = model.FamilyItem{
			ID:         core.Ptr(v.ID),
			Name:       v.Name,
			Status:     model.ShoppingStatus(v.Status),
			CategoryID: v.CategoryID,
			CreatedBy:  model.AccountID(v.CreatedBy),
			FamilyID:   model.FamilyID(v.FamilyID),
		}
	}

	return items
}

func NewFamilyItem(c *db.Conn) repository.FamilyItem {
	return &familyItemDao{conn: c}
}

func (dto *familyItemDto) TableName() string {
	return "family_items"
}
