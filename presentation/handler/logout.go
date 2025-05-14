package handler

import (
	"net/http"
	"sharebasket/presentation/response"
	"time"
)

func NewLogout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// アクセストークンのクッキーを削除
		accessCookie, err := r.Cookie("access_token")
		if err == nil && accessCookie != nil {
			http.SetCookie(w, &http.Cookie{
				Name:     "access_token",
				Value:    "",
				Path:     "/",
				MaxAge:   -1,
				HttpOnly: true,
				Secure:   r.TLS != nil,
				SameSite: http.SameSiteStrictMode,
				Expires:  time.Now().Add(-1 * time.Hour),
			})
		}

		// リフレッシュトークンのクッキーを削除
		refreshCookie, err := r.Cookie("refresh_token")
		if err == nil && refreshCookie != nil {
			http.SetCookie(w, &http.Cookie{
				Name:     "refresh_token",
				Value:    "",
				Path:     "/",
				MaxAge:   -1,
				HttpOnly: true,
				Secure:   r.TLS != nil,
				SameSite: http.SameSiteStrictMode,
				Expires:  time.Now().Add(-1 * time.Hour),
			})
		}

		response.NoContent(w)
	}
}
