package model

import (
	"errors"
	"fmt"
	"sharebasket/core"
	"slices"
)

const (
	NotPremiumMaxMember = 2
	PremiumMaxMember    = 4
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

// 家族にメンバーが存在するかどうかを判定。
func (f *Family) HasMembers() bool {
	return len(f.MemberIDs) != 0
}

// 新しいメンバーを家族に追加。
// オーナーのサブスクリプションタイプに応じたメンバー数の制限を超える場合にはエラー返す。
func (f *Family) Join(id AccountID) error {
	// すでに参加しているメンバーではないかチェック
	if slices.Contains(f.MemberIDs, id) {
		return core.NewInvalidError(errors.New("this account is already a member")).
			WithMessage("すでに参加しています")
	}

	// 家族に参加可能かチェック
	if !f.CanInvite() {
		return core.NewInvalidError(
			fmt.Errorf("cannnot join member: limit is %d", f.maxMembers()),
		).WithMessage("これ以上メンバーを追加できません")
	}

	f.MemberIDs = append(f.MemberIDs, id)
	return nil
}

// 家族に招待可能か判定する
func (f *Family) CanInvite() bool {
	return len(f.MemberIDs)+1 > f.maxMembers()
}

func (f *Family) maxMembers() int {
	if f.Owner.IsPremium {
		return PremiumMaxMember
	}
	return NotPremiumMaxMember
}
