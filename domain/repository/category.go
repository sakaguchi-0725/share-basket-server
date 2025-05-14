package repository

import "sharebasket/domain/model"

type Category interface {
	GetAll() ([]model.Category, error)
}
