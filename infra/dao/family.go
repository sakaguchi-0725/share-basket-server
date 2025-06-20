package dao

import (
	"context"
	"errors"
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

func (f *familyDao) HasMembership(ctx context.Context, accountID model.AccountID, familyID model.FamilyID) (bool, error) {
	var count int64
	err := f.conn.Model(&familyMemberDto{}).
		Where("account_id = ? AND family_id = ?", accountID.String(), familyID.String()).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (f *familyDao) GetByAccountID(ctx context.Context, id model.AccountID) (model.Family, error) {
	// オーナーかどうか確認
	isOwner, err := f.HasOwnedFamily(ctx, id)
	if err != nil {
		return model.Family{}, err
	}

	// オーナーの場合
	if isOwner {
		return f.GetOwnedFamily(ctx, id)
	}

	// メンバーの場合
	var member familyMemberDto
	err = f.conn.Preload("Family").
		Where("account_id = ?", id.String()).
		First(&member).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.Family{}, core.NewInvalidError(ErrRecordNotFound)
		}
		return model.Family{}, err
	}

	// 家族のメンバー一覧を取得
	var members []familyMemberDto
	err = f.conn.Where("family_id = ?", member.Family.ID).Find(&members).Error
	if err != nil {
		return model.Family{}, err
	}

	return member.Family.ToModel(members), nil
}

func (f *familyDao) GetByToken(ctx context.Context, token string) (model.Family, error) {
	id := f.client.Get(ctx, token).String()

	var family familyDto

	err := f.conn.Where("id = ?", id).First(&family).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.Family{}, core.NewInvalidError(ErrRecordNotFound)
		}

		return model.Family{}, err
	}

	var members []familyMemberDto

	err = f.conn.Where("family_id", family.ID).Find(&members).Error
	if err != nil {
		return model.Family{}, nil
	}

	return family.ToModel(members), nil
}

func (f *familyDao) HasOwnedFamily(ctx context.Context, accountID model.AccountID) (bool, error) {
	var count int64
	err := f.conn.Model(&familyDto{}).
		Where("owner_id = ?", accountID.String()).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (f *familyDao) GetOwnedFamily(ctx context.Context, accountID model.AccountID) (model.Family, error) {
	var family familyDto

	err := f.conn.Where("owner_id = ?", accountID.String()).First(&family).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.Family{}, core.NewInvalidError(ErrRecordNotFound)
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

// 指定されたアカウントIDが家族のオーナーまたはメンバーとして存在するかを確認する
func (f *familyDao) HasFamily(ctx context.Context, accountID model.AccountID) (bool, error) {
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

func (f *familyDao) Store(ctx context.Context, family *model.Family) error {
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
