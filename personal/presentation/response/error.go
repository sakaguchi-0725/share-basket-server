package response

import (
	"net/http"
	"share-basket-server/core/apperr"
)

type errorResponse struct {
	Code   string `json:"code"`
	Detail string `json:"detail"`
}

func Error(w http.ResponseWriter, err error) error {
	if appErr, ok := err.(*apperr.AppError); ok {
		var status int

		switch appErr.Code() {
		case apperr.ErrBadRequest:
		case apperr.ErrExpiredCode:
			status = http.StatusBadRequest
		case apperr.ErrUnauthorized:
			status = http.StatusUnauthorized
		case apperr.ErrNotFound:
			status = http.StatusNotFound
		default:
			status = http.StatusInternalServerError
		}

		return httpJSON(w, status, &errorResponse{
			Code:   appErr.Code().String(),
			Detail: appErr.Error(),
		})
	}

	return httpJSON(w, http.StatusInternalServerError, errorResponse{
		Code:   "InternalServerError",
		Detail: err.Error(),
	})
}
