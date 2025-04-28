package input

import (
	"context"
	"share-basket-server/personal/usecase/output"
)

type LoginInput struct {
	Email    string
	Password string
}

type LoginInputPort interface {
	Execute(ctx context.Context, in LoginInput, out output.LoginOutputPort) error
}
