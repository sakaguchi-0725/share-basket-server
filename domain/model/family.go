package model

import (
	"errors"
	"fmt"
	"slices"
)

const (
	NotPremiumMaxMember = 2
	PremiumMaxMember    = 4
)

var (
	ErrNonPremiumMaxMembers = errors.New("non-premium users can have up to 2 family members")
	ErrPremiumMaxMembers    = errors.New("premium users can have up to 4 family members")
)

type Family struct {
	ID        FamilyID
	Name      string
	Owner     Account
	MemberIDs []AccountID
}

// NewFamily は、指定されたID、名前、オーナーを持つ新しいFamilyインスタンスを作成します。
// 家族名が空である場合にエラーを返します。
func NewFamily(id FamilyID, name string, owner Account) (Family, error) {
	if name == "" {
		return Family{}, errors.New("family name is required")
	}

	return Family{
		ID:        id,
		Name:      name,
		Owner:     owner,
		MemberIDs: []AccountID{}, // 家族登録時はメンバー追加を行わない
	}, nil
}

// HasMembers は、家族にメンバーが存在するかどうかを判定します。
// メンバーが1人以上存在する場合はtrueを返します。
func (f *Family) HasMembers() bool {
	return len(f.MemberIDs) != 0
}

// Join は、新しいメンバーを家族に追加します。
// オーナーのサブスクリプションタイプに応じたメンバー数の制限を超える場合にはエラーを返します。
func (f *Family) Join(id AccountID) error {
	// すでに参加しているメンバーではないかチェック
	if slices.Contains(f.MemberIDs, id) {
		return errors.New("this account is already a member")
	}

	// 家族のメンバー上限を超えていないかチェック
	if len(f.MemberIDs)+1 > f.maxMembers() {
		return fmt.Errorf("cannnot join member: limit is %d", f.maxMembers())
	}

	f.MemberIDs = append(f.MemberIDs, id)
	return nil
}

func (f *Family) maxMembers() int {
	if f.Owner.IsPremium {
		return PremiumMaxMember
	}
	return NotPremiumMaxMember
}
