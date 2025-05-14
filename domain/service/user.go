package service

import "sharebasket/domain/repository"

type User interface {
	IsEmailAvailable(email string) (bool, error)
}

type userService struct {
	repo repository.User
}

// 入力されたemailがすでに登録されていないか検証する。
func (u *userService) IsEmailAvailable(email string) (bool, error) {
	exists, err := u.repo.ExistsByEmail(email)
	if err != nil {
		return false, err
	}

	return !exists, nil
}

func NewUser(r repository.User) User {
	return &userService{repo: r}
}
