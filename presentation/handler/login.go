package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"share-basket-server/core/apperr"
	"share-basket-server/presentation/response"
	"share-basket-server/presentation/validator"
	"share-basket-server/usecase"
)

type (
	loginRequest struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}

	loginPresenter struct {
		w http.ResponseWriter
	}
)

func MakeLoginHandler(
	usecase usecase.LoginInputPort,
	validator validator.RequestValidator,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req loginRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Error(w, apperr.NewInvalidError(err))
			return
		}

		if err := validator.Validate(&req); err != nil {
			response.Error(w, apperr.NewInvalidError(err))
			return
		}

		err := usecase.Execute(r.Context(), req.makeInput(), NewLoginPresenter(w))
		if err != nil {
			response.Error(w, err)
		}
	}
}

func (req loginRequest) makeInput() usecase.LoginInput {
	return usecase.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	}
}

func (presenter *loginPresenter) Render(ctx context.Context, token string) error {
	cookie := &http.Cookie{
		Name:     "access_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // dev環境ではfalse
		SameSite: http.SameSiteStrictMode,
		MaxAge:   86400, // 24時間（1日）
	}

	http.SetCookie(presenter.w, cookie)
	response.NoContent(presenter.w)

	return nil
}

func NewLoginPresenter(w http.ResponseWriter) usecase.LoginOutputPort {
	return &loginPresenter{w}
}
