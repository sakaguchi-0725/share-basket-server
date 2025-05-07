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
	signUpRequest struct {
		Name     string `json:"name" validate:"required"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}

	signUpPresenter struct {
		w http.ResponseWriter
	}
)

func MakeSignUpHandler(
	usecase usecase.SignUpInputPort,
	validator validator.RequestValidator,
	logger logger.Logger,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req signUpRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			logger.
				With("request body", r.Body).
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

		err := usecase.Execute(r.Context(), req.makeInput(), NewSignUpPresenter(w))
		if err != nil {
			response.Error(w, err)
		}
	}
}

func (req signUpRequest) makeInput() usecase.SignUpInput {
	return usecase.SignUpInput{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}
}

func (presenter *signUpPresenter) Render(ctx context.Context) error {
	response.NoContent(presenter.w)
	return nil
}

func NewSignUpPresenter(w http.ResponseWriter) usecase.SignUpConfirmOutputPort {
	return &signUpPresenter{w}
}
