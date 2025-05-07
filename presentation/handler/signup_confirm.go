package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"share-basket-server/core/apperr"
	"share-basket-server/core/logger"
	"share-basket-server/presentation/response"
	"share-basket-server/presentation/validator"
	"share-basket-server/usecase"
)

type (
	signUpConfirmRequest struct {
		Email            string `json:"email" validate:"required,email"`
		ConfirmationCode string `json:"confirmationCode" validate:"required"`
	}

	signUpConfirmPresenter struct {
		w http.ResponseWriter
	}
)

func MakeSignUpConfirmHandler(
	usecase usecase.SignUpConfirmInputPort,
	validator validator.RequestValidator,
	logger logger.Logger,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req signUpConfirmRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			logger.
				With("resuest body", r.Body).
				With("error", err).
				Info("invalid request")
			response.Error(w, apperr.NewInvalidError(err))
			return
		}

		if err := validator.Validate(&req); err != nil {
			logger.
				With("request body", req).
				With("error", err).
				Info("validation error")
			response.Error(w, apperr.NewInvalidError(err))
			return
		}

		err := usecase.Execute(r.Context(), req.makeInput(), NewSignUpConfirmPresenter(w))
		if err != nil {
			response.Error(w, err)
		}
	}
}

func (req signUpConfirmRequest) makeInput() usecase.SignUpConfirmInput {
	return usecase.SignUpConfirmInput{
		Email:            req.Email,
		ConfirmationCode: req.ConfirmationCode,
	}
}

func (presenter *signUpConfirmPresenter) Render(ctx context.Context) error {
	response.NoContent(presenter.w)
	return nil
}

func NewSignUpConfirmPresenter(w http.ResponseWriter) usecase.SignUpConfirmOutputPort {
	return &signUpConfirmPresenter{w}
}
