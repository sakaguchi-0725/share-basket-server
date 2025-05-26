package service

import (
	"context"
	"sharebasket/domain/repository"
)

type User interface {
	IsEmailAvailable(ctx context.Context, email string) (bool, error)
}

type userService struct {
	repo repository.User
}

// 入力されたemailがすでに登録されていないか検証する。
func (u *userService) IsEmailAvailable(ctx context.Context, email string) (bool, error) {
	exists, err := u.repo.ExistsByEmail(ctx, email)
	if err != nil {
		return false, err
	}

	return !exists, nil
}

func NewUser(r repository.User) User {
	return &userService{repo: r}
}
