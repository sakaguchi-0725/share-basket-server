package input

import (
	"context"
	"share-basket-server/personal/usecase/output"
)

type SignUpConfirmInput struct {
	Email            string
	ConfirmationCode string
}

type SignUpConfirmInputPort interface {
	Execute(ctx context.Context, in SignUpConfirmInput, out output.SignUpConfirmOutputPort) error
}
