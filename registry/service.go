package registry

import "sharebasket/domain/service"

type Service struct {
	User   service.User
	Family service.Family
}

func NewService(r *Repository) *Service {
	return &Service{
		User:   service.NewUser(r.User),
		Family: service.NewFamily(r.Family),
	}
}
