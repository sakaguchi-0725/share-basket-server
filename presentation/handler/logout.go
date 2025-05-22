package handler

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func NewLogout() echo.HandlerFunc {
	return func(c echo.Context) error {
		// アクセストークンのクッキーを削除
		accessCookie, err := c.Cookie("access_token")
		if err == nil && accessCookie != nil {
			c.SetCookie(&http.Cookie{
				Name:     "access_token",
				Value:    "",
				Path:     "/",
				MaxAge:   -1,
				HttpOnly: true,
				Secure:   c.Request().TLS != nil,
				SameSite: http.SameSiteStrictMode,
				Expires:  time.Now().Add(-1 * time.Hour),
			})
		}

		// リフレッシュトークンのクッキーを削除
		refreshCookie, err := c.Cookie("refresh_token")
		if err == nil && refreshCookie != nil {
			c.SetCookie(&http.Cookie{
				Name:     "refresh_token",
				Value:    "",
				Path:     "/",
				MaxAge:   -1,
				HttpOnly: true,
				Secure:   c.Request().TLS != nil,
				SameSite: http.SameSiteStrictMode,
				Expires:  time.Now().Add(-1 * time.Hour),
			})
		}

		return c.NoContent(http.StatusNoContent)
	}
}
