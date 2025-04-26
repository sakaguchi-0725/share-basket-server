package input

import (
	"context"
	"share-basket-server/personal/usecase/output"
)

type SignUpInput struct {
	Name     string
	Email    string
	Password string
}

type SignUpInputPort interface {
	Execute(ctx context.Context, in SignUpInput, out output.SignUpOutputPort) error
}
