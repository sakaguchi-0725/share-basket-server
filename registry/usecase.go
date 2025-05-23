package registry

import (
	"sharebasket/usecase"
)

type (
	UseCase interface {
		NewLogin() usecase.Login
		NewSignUp() usecase.SignUp
		NewSignUpConfirm() usecase.SignUpConfirm
		NewVerifyToken() usecase.VerifyToken
		NewGetAccount() usecase.GetAccount
		NewGetCategories() usecase.GetCategories
		NewCreatePersonalItem() usecase.CreatePersonalItem
		NewGetPersonalItems() usecase.GetPersonalItems
		NewUpdatePersonalItem() usecase.UpdatePersonalItem
		NewDeletePersonalItem() usecase.DeletePersonalItem
		NewCreateFamily() usecase.CreateFamily
		NewInvitationFamily() usecase.InvitationFamily
		NewJoinFamily() usecase.JoinFamily
	}

	usecaseImpl struct {
		repo    Repository
		service Service
	}
)

func (u *usecaseImpl) NewJoinFamily() usecase.JoinFamily {
	return usecase.NewJoinFamily(u.repo.NewAccount(), u.repo.NewFamily())
}

func (u *usecaseImpl) NewInvitationFamily() usecase.InvitationFamily {
	return usecase.NewInvitationFamily(
		u.repo.NewAccount(),
		u.repo.NewFamily(),
		u.service.NewFamily(),
	)
}

func (u *usecaseImpl) NewCreateFamily() usecase.CreateFamily {
	return usecase.NewCreateFamily(
		u.repo.NewAccount(),
		u.repo.NewFamily(),
		u.service.NewFamily(),
	)
}

func (u *usecaseImpl) NewDeletePersonalItem() usecase.DeletePersonalItem {
	return usecase.NewDeletePersonalItem(
		u.repo.NewAccount(),
		u.repo.NewPersonalItem(),
	)
}

func (u *usecaseImpl) NewUpdatePersonalItem() usecase.UpdatePersonalItem {
	return usecase.NewUpdatePersonalItem(
		u.repo.NewAccount(),
		u.repo.NewPersonalItem(),
	)
}

func (u *usecaseImpl) NewGetPersonalItems() usecase.GetPersonalItems {
	return usecase.NewGetPersonalItems(
		u.repo.NewAccount(),
		u.repo.NewPersonalItem(),
	)
}

func (u *usecaseImpl) NewCreatePersonalItem() usecase.CreatePersonalItem {
	return usecase.NewCreatePersonalItem(
		u.repo.NewAccount(),
		u.repo.NewPersonalItem(),
	)
}

func (u *usecaseImpl) NewGetAccount() usecase.GetAccount {
	return usecase.NewGetAccount(u.repo.NewAccount())
}

func (u *usecaseImpl) NewGetCategories() usecase.GetCategories {
	return usecase.NewGetCategories(u.repo.NewCategory())
}

func (u *usecaseImpl) NewLogin() usecase.Login {
	return usecase.NewLogin(u.repo.NewAuthenticator())
}

func (u *usecaseImpl) NewSignUp() usecase.SignUp {
	return usecase.NewSignUp(
		u.repo.NewAuthenticator(),
		u.repo.NewUser(),
		u.repo.NewAccount(),
		u.service.NewUser(),
		u.repo.NewTransaction(),
	)
}

func (u *usecaseImpl) NewSignUpConfirm() usecase.SignUpConfirm {
	return usecase.NewSignUpConfirm(u.repo.NewAuthenticator())
}

func (u *usecaseImpl) NewVerifyToken() usecase.VerifyToken {
	return usecase.NewVerifyToken(u.repo.NewAuthenticator())
}

func NewUseCase(r Repository, s Service) UseCase {
	return &usecaseImpl{
		repo:    r,
		service: s,
	}
}
