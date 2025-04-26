package usecase

import (
	"context"
	"errors"
	"share-basket-server/core/apperr"
	"share-basket-server/personal/domain/repository"
	"share-basket-server/personal/usecase/input"
	"share-basket-server/personal/usecase/output"
)

type signUpConfirmUseCase struct {
	authenticator repository.Authenticator
}

func (s *signUpConfirmUseCase) Execute(
	ctx context.Context, in input.SignUpConfirmInput, out output.SignUpConfirmOutputPort,
) error {
	err := s.authenticator.SignUpConfirm(ctx, in.Email, in.ConfirmationCode)

	if err != nil {
		if errors.Is(err, apperr.ErrInvalidData) {
			return apperr.NewInvalidError(err)
		}
		if errors.Is(err, apperr.ErrExpiredCodeException) {
			return apperr.New(apperr.ErrExpiredCode, err)
		}
		return err
	}

	return out.Render(ctx)
}

func NewSignUpConfirmUseCase(authenticator repository.Authenticator) input.SignUpConfirmInputPort {
	return &signUpConfirmUseCase{authenticator}
}
