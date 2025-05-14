package dao

import (
	"sharebasket/domain/model"
	"sharebasket/domain/repository"
	"sharebasket/infra/db"
	"time"
)

type (
	categoryDto struct {
		ID        int64     `gorm:"primaryKey autoIncrement"`
		Name      string    `gorm:"not null"`
		CreatedAt time.Time `gorm:"autoCreateTime"`
		UpdatedAt time.Time `gorm:"autoUpdateTime"`
	}

	categoryDtos []categoryDto

	categoryDao struct {
		conn *db.Conn
	}
)

func (c *categoryDao) GetAll() ([]model.Category, error) {
	var dtos categoryDtos

	err := c.conn.Find(&dtos).Error
	if err != nil {
		return []model.Category{}, err
	}

	return dtos.toModels(), nil
}

func (dtos categoryDtos) toModels() []model.Category {
	categories := make([]model.Category, len(dtos))

	for i, v := range dtos {
		categories[i] = model.Category{
			ID:   v.ID,
			Name: v.Name,
		}
	}

	return categories
}

func NewCategory(c *db.Conn) repository.Category {
	return &categoryDao{conn: c}
}

func (c *categoryDto) TableName() string {
	return "categories"
}
