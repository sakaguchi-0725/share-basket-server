package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"share-basket-server/core/apperr"
	"share-basket-server/personal/presentation/response"
	"share-basket-server/personal/usecase"
)

type (
	signUpConfirmRequest struct {
		Email            string `json:"email"`
		ConfirmationCode string `json:"confirmationCode"`
	}

	signUpConfirmPresenter struct {
		w http.ResponseWriter
	}
)

func MakeSignUpConfirmHandler(usecase usecase.SignUpConfirmInputPort) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req signUpConfirmRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
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
