package middleware

import (
	"context"
	"net/http"
	"share-basket-server/core/apperr"
	contextKey "share-basket-server/core/context"
	"share-basket-server/presentation/response"
	"share-basket-server/usecase"
)

func Auth(verifyToken usecase.VerifyTokenInputPort) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("access_token")
			if err != nil {
				if err == http.ErrNoCookie {
					response.Error(w, apperr.New(apperr.ErrUnauthorized, err))
					return
				}
				response.Error(w, err)
				return
			}

			userID, err := verifyToken.Execute(r.Context(), cookie.Value)
			if err != nil {
				response.Error(w, err)
				return
			}

			ctx := context.WithValue(r.Context(), contextKey.UserID, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
