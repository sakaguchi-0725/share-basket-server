package registry

import (
	"context"
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
	}

	repositoryImpl struct {
		conn   *db.Conn
		client auth.CognitoClient
	}
)

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

func NewRepository(c *db.Conn) (Repository, error) {
	client, err := auth.NewCognitoClient(context.Background())
	if err != nil {
		return nil, err
	}

	return &repositoryImpl{
		conn:   c,
		client: client,
	}, nil
}
