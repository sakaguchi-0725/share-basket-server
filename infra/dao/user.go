package dao

import (
	"context"
	"log"
	"sharebasket/domain/model"
	"sharebasket/domain/repository"
	"sharebasket/infra/db"
	"sharebasket/infra/transaction"
	"time"
)

type (
	userDto struct {
		ID        string    `gorm:"primaryKey"`
		Email     string    `gorm:"not null"`
		CreatedAt time.Time `gorm:"autoCreateTime"`
		UpdatedAt time.Time `gorm:"autoUpdateTime"`
	}

	user struct {
		conn *db.Conn
	}
)

// 指定されたメールアドレスを持つユーザーが存在するかを確認する。
func (u *user) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64

	err := u.conn.Model(&userDto{}).Where("email = ?", email).Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// Store はユーザー情報を保存する。
// 新規登録と更新の両方に対応しており、IDの有無で判断する。
func (u *user) Store(ctx context.Context, user *model.User) error {
	dto := newUserDto(user)
	tx := transaction.GetTx(ctx)

	if err := tx.Save(&dto).Error; err != nil {
		return err
	}

	return nil
}

// StoreWithTx はトランザクション内でユーザー情報を保存する。
func (u *user) StoreWithTx(tx *db.Conn, user *model.User) error {
	dto := newUserDto(user)
	log.Println(dto)

	if err := tx.Save(&dto).Error; err != nil {
		return err
	}

	return nil
}

func NewUser(c *db.Conn) repository.User {
	return &user{conn: c}
}

func newUserDto(u *model.User) userDto {
	return userDto{
		ID:    u.ID,
		Email: u.Email,
	}
}

func (u userDto) TableName() string {
	return "users"
}
