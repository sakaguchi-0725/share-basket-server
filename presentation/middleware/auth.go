package middleware

import (
	"context"
	"net/http"
	"sharebasket/core"
	"sharebasket/presentation/response"
	"sharebasket/usecase"
)

func Auth(v usecase.VerifyToken, logger core.Logger) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("access_token")
			if err != nil {
				if err == http.ErrNoCookie {
					logger.WithError(err).Warn("access token is not found")
					response.Error(w, core.NewAppError(core.ErrUnauthorized, err))
					return
				}

				logger.WithError(err).Error("failed to get cookie")
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
