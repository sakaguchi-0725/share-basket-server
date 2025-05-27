package registry

import "sharebasket/domain/service"

type Service struct {
	User    service.User
	Account service.Account
}

func NewService(r *Repository) *Service {
	return &Service{
		User:    service.NewUser(r.User),
		Account: service.NewAccount(r.Family),
	}
}
