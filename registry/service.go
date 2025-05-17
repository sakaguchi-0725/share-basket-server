package registry

import "sharebasket/domain/service"

type (
	Service interface {
		NewUser() service.User
		NewFamily() service.Family
	}

	serviceImpl struct {
		repo Repository
	}
)

func (s *serviceImpl) NewFamily() service.Family {
	return service.NewFamily(s.repo.NewFamily())
}

func (s *serviceImpl) NewUser() service.User {
	return service.NewUser(s.repo.NewUser())
}

func NewService(r Repository) Service {
	return &serviceImpl{repo: r}
}
