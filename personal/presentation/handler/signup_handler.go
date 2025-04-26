package handler

import (
	"encoding/json"
	"net/http"
	"share-basket-server/core/middleware"
	"share-basket-server/personal/presentation/presenter"
	"share-basket-server/personal/usecase/input"
)

type SignUpRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func MakeSignUpHandler(usecase input.SignUpInputPort) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req SignUpRequest
		errRw := w.(*middleware.ErrResponseWriter)

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			errRw.Err = err
			return
		}

		out := presenter.NewSignUpPresenter(w)
		ctx := r.Context()

		errRw.Err = usecase.Execute(ctx, req.makeInput(), out)
	}
}

func (req SignUpRequest) makeInput() input.SignUpInput {
	return input.SignUpInput{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}
}
