package middleware

import (
	"encoding/json"
	"errors"
	"net/http"
	"share-basket-server/core/apperr"
	"share-basket-server/personal/presentation/response"
)

type ErrorResponse struct {
	Code string `json:"code"`
}

func WithError(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ew := &response.ErrResponseWriter{ResponseWriter: w}
		next.ServeHTTP(ew, r)

		if ew.Err == nil {
			return
		}

		w.Header().Set("Content-Type", "application/json")

		var appErr *apperr.AppError
		if errors.As(ew.Err, &appErr) {
			switch appErr.Code() {
			case apperr.ErrBadRequest.String():
				w.WriteHeader(http.StatusBadRequest)
			case apperr.ErrNotFound.String():
				w.WriteHeader(http.StatusNotFound)
			case apperr.ErrUnauthorized.String():
				w.WriteHeader(http.StatusUnauthorized)
			default:
				w.WriteHeader(http.StatusInternalServerError)
			}

			json.NewEncoder(w).Encode(ErrorResponse{
				Code: appErr.Code(),
			})
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{
			Code: "InternalServer",
		})
	})
}
