package dao

import (
	"context"
	"errors"
	"sharebasket/domain/model"
	"sharebasket/domain/repository"
	"sharebasket/infra/db"
	"sharebasket/infra/transaction"
	"sharebasket/usecase"
	"time"

	"gorm.io/gorm"
)

type (
	accountDto struct {
		ID        string    `gorm:"primaryKey"`
		Name      string    `gorm:"not null"`
		UserID    string    `gorm:"not null"`
		User      userDto   `gorm:"foreignKey:UserID;references:ID"`
		CreatedAt time.Time `gorm:"autoCreateTime"`
		UpdatedAt time.Time `gorm:"autoUpdateTime"`
	}

	account struct {
		conn *db.Conn
	}
)

func (a *account) Get(ctx context.Context, userID string) (model.Account, error) {
	var dto accountDto

	err := a.conn.Where("user_id = ?", userID).First(&dto).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.Account{}, usecase.ErrAccountNotFound
		}
		return model.Account{}, err
	}

	return dto.toModel(), nil
}

func (a *account) Store(ctx context.Context, acc *model.Account) error {
	dto := newAccountDto(acc)
	tx := transaction.GetTx(ctx)

	err := tx.Save(&dto).Error
	if err != nil {
		return err
	}

	return nil
}

func NewAccount(c *db.Conn) repository.Account {
	return &account{conn: c}
}

func newAccountDto(a *model.Account) accountDto {
	return accountDto{
		ID:     a.ID.String(),
		Name:   a.Name,
		UserID: a.UserID,
	}
}

func (dto accountDto) toModel() model.Account {
	return model.Account{
		ID:     model.AccountID(dto.ID),
		Name:   dto.Name,
		UserID: dto.UserID,
	}
}

func (a accountDto) TableName() string {
	return "accounts"
}
