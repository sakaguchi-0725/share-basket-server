package registry

import (
	"sharebasket/usecase"
)

type UseCase struct {
	Login              usecase.Login
	SignUp             usecase.SignUp
	SignUpConfirm      usecase.SignUpConfirm
	VerifyToken        usecase.VerifyToken
	GetAccount         usecase.GetAccount
	GetCategories      usecase.GetCategories
	CreatePersonalItem usecase.CreatePersonalItem
	GetPersonalItems   usecase.GetPersonalItems
	UpdatePersonalItem usecase.UpdatePersonalItem
	DeletePersonalItem usecase.DeletePersonalItem
	CreateFamily       usecase.CreateFamily
	InvitationFamily   usecase.InvitationFamily
	JoinFamily         usecase.JoinFamily
	CreateFamilyItem   usecase.CreateFamilyItem
}

func NewUseCase(r *Repository, s *Service) *UseCase {
	return &UseCase{
		Login: usecase.NewLogin(r.Authenticator),
		SignUp: usecase.NewSignUp(
			r.Authenticator,
			r.User,
			r.Account,
			s.User,
			r.Transaction,
		),
		SignUpConfirm: usecase.NewSignUpConfirm(r.Authenticator),
		VerifyToken:   usecase.NewVerifyToken(r.Authenticator),
		GetAccount:    usecase.NewGetAccount(r.Account),
		GetCategories: usecase.NewGetCategories(r.Category),
		CreatePersonalItem: usecase.NewCreatePersonalItem(
			r.Account,
			r.PersonalItem,
		),
		GetPersonalItems: usecase.NewGetPersonalItems(
			r.Account,
			r.PersonalItem,
		),
		UpdatePersonalItem: usecase.NewUpdatePersonalItem(
			r.Account,
			r.PersonalItem,
		),
		DeletePersonalItem: usecase.NewDeletePersonalItem(
			r.Account,
			r.PersonalItem,
		),
		CreateFamily: usecase.NewCreateFamily(
			r.Account,
			r.Family,
			s.Family,
		),
		InvitationFamily: usecase.NewInvitationFamily(
			r.Account,
			r.Family,
			s.Family,
		),
		JoinFamily: usecase.NewJoinFamily(r.Account, r.Family),
		CreateFamilyItem: usecase.NewCreateFamilyItem(
			r.Account,
			r.Family,
			r.FamilyItem,
			s.Family,
		),
	}
}
