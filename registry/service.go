package registry

import "sharebasket/domain/service"

type (
	Service interface {
		NewUser() service.User
	}

	serviceImpl struct {
		repo Repository
	}
)

func (s *serviceImpl) NewUser() service.User {
	return service.NewUser(s.repo.NewUser())
}

func NewService(r Repository) Service {
	return &serviceImpl{repo: r}
}
