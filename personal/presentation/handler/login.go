package handler

import (
	"encoding/json"
	"net/http"
	"share-basket-server/core/apperr"
	"share-basket-server/personal/presentation/presenter"
	"share-basket-server/personal/presentation/response"
	"share-basket-server/personal/usecase/input"
)

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func MakeLoginHandler(usecase input.LoginInputPort) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req loginRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Error(w, apperr.NewInvalidError(err))
			return
		}

		out := presenter.NewLoginPresenter(w)
		ctx := r.Context()

		err := usecase.Execute(ctx, req.makeInput(), out)
		if err != nil {
			response.Error(w, err)
		}
	}
}

func (req loginRequest) makeInput() input.LoginInput {
	return input.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	}
}
