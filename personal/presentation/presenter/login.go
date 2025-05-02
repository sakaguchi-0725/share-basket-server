package presenter

import (
	"context"
	"net/http"
	"share-basket-server/personal/presentation/response"
	"share-basket-server/personal/usecase/output"
)

type loginPresenter struct {
	w http.ResponseWriter
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

func NewLoginPresenter(w http.ResponseWriter) output.LoginOutputPort {
	return &loginPresenter{w}
}
