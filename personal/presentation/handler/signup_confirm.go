package handler

import (
	"encoding/json"
	"net/http"
	"share-basket-server/core/apperr"
	"share-basket-server/personal/presentation/presenter"
	"share-basket-server/personal/presentation/response"
	"share-basket-server/personal/usecase/input"
)

type signUpConfirmRequest struct {
	Email            string `json:"email"`
	ConfirmationCode string `json:"confirmationCode"`
}

func MakeSignUpConfirmHandler(usecase input.SignUpConfirmInputPort) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req signUpConfirmRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Error(w, apperr.NewInvalidError(err))
			return
		}

		out := presenter.NewSignUpConfirmPresenter(w)
		ctx := r.Context()

		err := usecase.Execute(ctx, req.makeInput(), out)
		if err != nil {
			response.Error(w, err)
		}
	}
}

func (req signUpConfirmRequest) makeInput() input.SignUpConfirmInput {
	return input.SignUpConfirmInput{
		Email:            req.Email,
		ConfirmationCode: req.ConfirmationCode,
	}
}
