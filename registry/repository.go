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

type (
	Repository interface {
		NewAuthenticator() repository.Authenticator
		NewAccount() repository.Account
		NewUser() repository.User
		NewCategory() repository.Category
		NewPersonalItem() repository.PersonalItem
		NewTransaction() repository.Transaction
		NewFamily() repository.Family
		NewFamilyItem() repository.FamilyItem
	}

	repositoryImpl struct {
		conn      *db.Conn
		client    auth.CognitoClient
		redisHost string
	}
)

func (r *repositoryImpl) NewFamilyItem() repository.FamilyItem {
	return dao.NewFamilyItem(r.conn)
}

func (r *repositoryImpl) NewFamily() repository.Family {
	return dao.NewFamily(r.conn, r.redisHost)
}

func (r *repositoryImpl) NewPersonalItem() repository.PersonalItem {
	return dao.NewPersonalItem(r.conn)
}

func (r *repositoryImpl) NewCategory() repository.Category {
	return dao.NewCategory(r.conn)
}

func (r *repositoryImpl) NewTransaction() repository.Transaction {
	return transaction.New(r.conn)
}

func (r *repositoryImpl) NewAccount() repository.Account {
	return dao.NewAccount(r.conn)
}

func (r *repositoryImpl) NewAuthenticator() repository.Authenticator {
	return auth.NewCognito(r.client)
}

func (r *repositoryImpl) NewUser() repository.User {
	return dao.NewUser(r.conn)
}

func NewRepository(c *db.Conn, cfg *core.Config) (Repository, error) {
	client, err := auth.NewCognitoClient(context.Background(), cfg.AWS)
	if err != nil {
		return nil, err
	}

	return &repositoryImpl{
		conn:      c,
		client:    client,
		redisHost: cfg.RedisHost,
	}, nil
}
