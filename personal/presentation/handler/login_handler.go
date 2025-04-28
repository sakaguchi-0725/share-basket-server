package handler

import (
	"encoding/json"
	"net/http"
	"share-basket-server/personal/presentation/presenter"
	"share-basket-server/personal/presentation/response"
	"share-basket-server/personal/usecase/input"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func MakeLoginHandler(usecase input.LoginInputPort) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req LoginRequest
		errRw := w.(*response.ErrResponseWriter)

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			errRw.Err = err
			return
		}

		out := presenter.NewLoginPresenter(w)
		ctx := r.Context()

		errRw.Err = usecase.Execute(ctx, req.makeInput(), out)
	}
}

func (req LoginRequest) makeInput() input.LoginInput {
	return input.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	}
}
