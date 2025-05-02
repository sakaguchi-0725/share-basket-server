package handler

import (
	"encoding/json"
	"net/http"
	"share-basket-server/core/apperr"
	"share-basket-server/personal/presentation/presenter"
	"share-basket-server/personal/presentation/response"
	"share-basket-server/personal/usecase/input"
)

type signUpRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func MakeSignUpHandler(usecase input.SignUpInputPort) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req signUpRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Error(w, apperr.NewInvalidError(err))
			return
		}

		out := presenter.NewSignUpPresenter(w)
		ctx := r.Context()

		err := usecase.Execute(ctx, req.makeInput(), out)
		if err != nil {
			response.Error(w, err)
		}
	}
}

func (req signUpRequest) makeInput() input.SignUpInput {
	return input.SignUpInput{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}
}
