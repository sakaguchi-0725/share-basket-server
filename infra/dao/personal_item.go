package dao

import (
	"errors"
	"sharebasket/core"
	"sharebasket/domain/model"
	"sharebasket/domain/repository"
	"sharebasket/infra/db"
	"time"

	"gorm.io/gorm"
)

type (
	personalItemDto struct {
		ID         int64       `gorm:"primaryKey autoIncrement"`
		Name       string      `gorm:"not null"`
		Status     string      `gorm:"not null"`
		CategoryID int64       `gorm:"not null"`
		Category   categoryDto `gorm:"foreignKey:CategoryID;references:ID"`
		AccountID  string      `gorm:"not null"`
		Account    accountDto  `gorm:"foreignKey:AccountID;references:ID"`
		CreatedAt  time.Time   `gorm:"autoCreateTime"`
		UpdatedAt  time.Time   `gorm:"autoUpdateTime"`
	}

	personalItemDtos []personalItemDto

	personalItemDao struct {
		conn *db.Conn
	}
)

func (p *personalItemDao) Delete(id int64) error {
	result := p.conn.Delete(&personalItemDto{}, id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return core.ErrDataNotFound
	}

	return nil
}

func (p *personalItemDao) GetByID(id int64) (model.PersonalItem, error) {
	var dto personalItemDto

	err := p.conn.Where("id = ?", id).First(&dto).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.PersonalItem{}, core.ErrDataNotFound
		}
		return model.PersonalItem{}, err
	}

	return dto.toModel(), nil
}

func (p *personalItemDao) GetAll(accID model.AccountID, status *model.ShoppingStatus) ([]model.PersonalItem, error) {
	var dtos personalItemDtos

	query := p.conn.DB
	if status != nil {
		query = query.Where("status = ?", status.String())
	}

	if err := query.Find(&dtos).Error; err != nil {
		return []model.PersonalItem{}, err
	}

	return dtos.toModels(), nil
}

func (p *personalItemDao) Store(item *model.PersonalItem) error {
	dto := newPersonalItemDto(item)

	if err := p.conn.Save(&dto).Error; err != nil {
		return err
	}

	if item.ID == nil {
		item.ID = core.Ptr(dto.ID)
	}

	return nil
}

func newPersonalItemDto(item *model.PersonalItem) personalItemDto {
	return personalItemDto{
		ID:         core.Derefer(item.ID),
		Name:       item.Name,
		Status:     item.Status.String(),
		CategoryID: item.CategoryID,
		AccountID:  item.AccountID.String(),
	}
}

func (dto personalItemDto) toModel() model.PersonalItem {
	return model.PersonalItem{
		ID:         core.Ptr(dto.ID),
		Name:       dto.Name,
		Status:     model.ShoppingStatus(dto.Status),
		CategoryID: dto.CategoryID,
		AccountID:  model.AccountID(dto.AccountID),
	}
}

func (dtos personalItemDtos) toModels() []model.PersonalItem {
	items := make([]model.PersonalItem, len(dtos))

	for i, v := range dtos {
		items[i] = model.PersonalItem{
			ID:         core.Ptr(v.ID),
			Name:       v.Name,
			Status:     model.ShoppingStatus(v.Status),
			CategoryID: v.CategoryID,
			AccountID:  model.AccountID(v.AccountID),
		}
	}

	return items
}

func NewPersonalItem(c *db.Conn) repository.PersonalItem {
	return &personalItemDao{conn: c}
}

func (p *personalItemDto) TableName() string {
	return "personal_shopping_items"
}
