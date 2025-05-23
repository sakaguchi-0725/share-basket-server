package handler

import (
	"net/http"
	"sharebasket/core"
	"sharebasket/usecase"

	"github.com/labstack/echo/v4"
)

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewLogin(usecase usecase.Login) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req loginRequest

		if err := c.Bind(&req); err != nil {
			return core.NewInvalidError(err)
		}

		output, err := usecase.Execute(c.Request().Context(), req.makeInput())
		if err != nil {
			return err
		}

		setCookies(c, output)
		return c.NoContent(http.StatusNoContent)
	}
}

func (l *loginRequest) makeInput() usecase.LoginInput {
	return usecase.LoginInput{
		Email:    l.Email,
		Password: l.Password,
	}
}

func setCookies(c echo.Context, o usecase.LoginOutput) {
	accessToken := &http.Cookie{
		Name:     "access_token",
		Value:    o.AccessToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   c.Request().TLS != nil,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   86400,
	}

	refreshToken := &http.Cookie{
		Name:     "refresh_token",
		Value:    o.RefreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   c.Request().TLS != nil,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   86400,
	}

	c.SetCookie(accessToken)
	c.SetCookie(refreshToken)
}
