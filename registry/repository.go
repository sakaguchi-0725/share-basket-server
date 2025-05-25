package registry

import (
	"context"
	"sharebasket/core"
	"sharebasket/domain/repository"
	"sharebasket/infra/auth"
	"sharebasket/infra/dao"
	"sharebasket/infra/db"
	"sharebasket/infra/transaction"
)

type Repository struct {
	Authenticator repository.Authenticator
	Account       repository.Account
	User          repository.User
	Category      repository.Category
	PersonalItem  repository.PersonalItem
	Transaction   repository.Transaction
	Family        repository.Family
	FamilyItem    repository.FamilyItem
}

func NewRepository(c *db.Conn, cfg *core.Config) (*Repository, error) {
	client, err := auth.NewCognitoClient(context.Background(), cfg.AWS)
	if err != nil {
		return nil, err
	}

	return &Repository{
		Authenticator: auth.NewCognito(client),
		Account:       dao.NewAccount(c),
		User:          dao.NewUser(c),
		Category:      dao.NewCategory(c),
		PersonalItem:  dao.NewPersonalItem(c),
		Transaction:   transaction.New(c),
		Family:        dao.NewFamily(c, cfg.RedisHost),
		FamilyItem:    dao.NewFamilyItem(c),
	}, nil
}
