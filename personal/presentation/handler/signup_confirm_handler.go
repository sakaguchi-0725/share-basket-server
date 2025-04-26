package handler

import (
	"encoding/json"
	"net/http"
	"share-basket-server/core/middleware"
	"share-basket-server/personal/presentation/presenter"
	"share-basket-server/personal/usecase/input"
)

type SignUpConfirmRequest struct {
	Email            string `json:"email"`
	ConfirmationCode string `json:"confirmation_code"`
}

func MakeSignUpConfirmHandler(usecase input.SignUpConfirmInputPort) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req SignUpConfirmRequest
		errRw := w.(*middleware.ErrResponseWriter)

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			errRw.Err = err
			return
		}

		out := presenter.NewSignUpConfirmPresenter(w)
		ctx := r.Context()

		errRw.Err = usecase.Execute(ctx, req.makeInput(), out)
	}
}

func (req SignUpConfirmRequest) makeInput() input.SignUpConfirmInput {
	return input.SignUpConfirmInput{
		Email:            req.Email,
		ConfirmationCode: req.ConfirmationCode,
	}
}
