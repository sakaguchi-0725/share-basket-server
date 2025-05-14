package middleware

import (
	"context"
	"net/http"
	"sharebasket/core"
	"sharebasket/presentation/response"
	"sharebasket/usecase"
)

func Auth(v usecase.VerifyToken) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("access_token")
			if err != nil {
				if err == http.ErrNoCookie {
					response.Error(w, core.NewAppError(core.ErrUnauthorized, err))
					return
				}
				response.Error(w, err)
				return
			}

			userID, err := v.Execute(r.Context(), cookie.Value)
			if err != nil {
				response.Error(w, err)
				return
			}

			ctx := context.WithValue(r.Context(), core.UserIDKey, userID)
			h.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
