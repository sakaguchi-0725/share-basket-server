package dao

import (
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
)

// Get implements repository.FamilyItem.
func (f *familyItemDao) Get(id model.FamilyID) ([]model.FamilyItem, error) {
	panic("unimplemented")
}

func (f *familyItemDao) Store(item *model.FamilyItem) error {
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

func NewFamilyItem(c *db.Conn) repository.FamilyItem {
	return &familyItemDao{conn: c}
}

func (dto *familyItemDto) TableName() string {
	return "family_items"
}
