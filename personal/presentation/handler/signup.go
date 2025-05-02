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
	signUpRequest struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	signUpPresenter struct {
		w http.ResponseWriter
	}
)

func MakeSignUpHandler(usecase usecase.SignUpInputPort) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req signUpRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
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
