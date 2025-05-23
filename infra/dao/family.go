package dao

import (
	"context"
	"errors"
	"fmt"
	"sharebasket/core"
	"sharebasket/domain/model"
	"sharebasket/domain/repository"
	"sharebasket/infra/db"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type (
	familyDto struct {
		ID        string    `gorm:"primaryKey"`
		Name      string    `gorm:"not null"`
		OwnerID   string    `gorm:"not null"`
		CreatedAt time.Time `gorm:"autoCreateTime"`
		UpdatedAt time.Time `gorm:"autoUpdateTime"`
	}

	familyMemberDto struct {
		ID        int64      `gorm:"primaryKey autoIncrement"`
		FamilyID  string     `gorm:"not null"`
		Family    familyDto  `gorm:"foreignKey:FamilyID;references:ID"`
		AccountID string     `gorm:"not null"`
		Account   accountDto `gorm:"foreignKey:AccountID;references:ID"`
		CreatedAt time.Time  `gorm:"autoCreateTime"`
		UpdatedAt time.Time  `gorm:"autoUpdateTime"`
	}

	familyDao struct {
		client   *redis.Client
		conn     *db.Conn
		duration time.Duration
	}
)

func (f *familyDao) HasOwnedFamily(accountID model.AccountID) (bool, error) {
	var count int64
	err := f.conn.Model(&familyDto{}).
		Where("owner_id = ?", accountID.String()).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (f *familyDao) GetOwnedFamily(accountID model.AccountID) (model.Family, error) {
	var family familyDto

	err := f.conn.Where("owner_id = ?", accountID.String()).First(&family).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.Family{}, core.ErrDataNotFound
		}
		return model.Family{}, err
	}

	var members []familyMemberDto
	err = f.conn.Where("family_id = ?", family.ID).Find(&members).Error
	if err != nil {
		return model.Family{}, err
	}

	return family.ToModel(members), nil
}

func (f *familyDao) Invitation(ctx context.Context, token string, familyID model.FamilyID) error {
	err := f.client.Set(ctx, token, familyID.String(), f.duration).Err()
	if err != nil {
		return err
	}
	return nil
}

// HasFamily は指定されたアカウントIDが家族のオーナーまたはメンバーとして存在するかを確認する
func (f *familyDao) HasFamily(accountID model.AccountID) (bool, error) {
	var count int64
	err := f.conn.Model(&familyDto{}).
		Joins("LEFT JOIN family_members ON families.id = family_members.family_id").
		Where("families.owner_id = ? OR family_members.account_id = ?", accountID.String(), accountID.String()).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (f *familyDao) Join(family *model.Family) error {
	var existsMembers []familyMemberDto

	err := f.conn.Where("family_id = ?", family.ID).Find(&existsMembers).Error
	if err != nil {
		return err
	}

	// 既存のメンバーIDを検索するためのマップ
	existingMemberIDs := make(map[string]bool)
	for _, member := range existsMembers {
		existingMemberIDs[member.AccountID] = true
	}

	// 追加が必要な新しいメンバーを保持するスライス
	var newMembers []familyMemberDto
	for _, memberID := range family.MemberIDs {
		// memberIDが既存のメンバーに含まれていないか確認
		if !existingMemberIDs[memberID.String()] {
			newMembers = append(newMembers, familyMemberDto{
				FamilyID:  family.ID.String(),
				AccountID: memberID.String(),
			})
		}
	}

	if len(newMembers) == 0 {
		// ドメイン側でバリデーション済みのため、通常は通らない。
		// 念のため明示的に分岐。
		return nil
	}

	if err := f.conn.Create(&newMembers).Error; err != nil {
		return fmt.Errorf("failed to insert new members: %w", err)
	}

	return nil
}

func (f *familyDao) Store(family *model.Family) error {
	dto := newFamilyDto(family)

	if err := f.conn.Save(&dto).Error; err != nil {
		return err
	}

	return nil
}

func newFamilyDto(f *model.Family) familyDto {
	return familyDto{
		ID:      f.ID.String(),
		Name:    f.Name,
		OwnerID: f.Owner.ID.String(),
	}
}

func newFamilyMemberDtos(f *model.Family) []familyMemberDto {
	members := make([]familyMemberDto, len(f.MemberIDs))

	for i, v := range f.MemberIDs {
		members[i] = familyMemberDto{
			FamilyID:  f.ID.String(),
			AccountID: v.String(),
		}
	}

	return members
}

func (f *familyDto) ToModel(members []familyMemberDto) model.Family {
	familyID := model.FamilyID(f.ID)
	ownerID := model.AccountID(f.OwnerID)
	owner := model.Account{ID: ownerID}

	memberIDs := make([]model.AccountID, len(members))
	for i, member := range members {
		memberIDs[i] = model.AccountID(member.AccountID)
	}

	return model.Family{
		ID:        familyID,
		Name:      f.Name,
		Owner:     owner,
		MemberIDs: memberIDs,
	}
}

func NewFamily(c *db.Conn, redisHost string) repository.Family {
	return &familyDao{
		client: redis.NewClient(&redis.Options{
			Addr:     redisHost,
			Password: "",
			DB:       0,
		}),
		conn:     c,
		duration: 6 * time.Hour,
	}
}

func (f familyDto) TableName() string {
	return "families"
}
