package response

import (
	"net/http"
	"sharebasket/core"
)

type errorResponse struct {
	Code   string `json:"code"`
	Detail string `json:"detail"`
}

func Error(w http.ResponseWriter, err error) error {
	if appErr, ok := err.(*core.AppError); ok {
		var status int

		switch appErr.Code() {
		case core.ErrBadRequest:
			status = http.StatusBadRequest
		case core.ErrEmailAlreadyExists:
			status = http.StatusConflict
		case core.ErrExpiredCode:
			status = http.StatusGone
		case core.ErrUnauthorized, core.ErrExpiredToken:
			status = http.StatusUnauthorized
		case core.ErrForbidden:
			status = http.StatusForbidden
		case core.ErrNotFound:
			status = http.StatusNotFound
		default:
			status = http.StatusInternalServerError
		}

		return httpJSON(w, status, &errorResponse{
			Code:   appErr.Code().String(),
			Detail: appErr.Error(),
		})
	}

	return httpJSON(w, http.StatusInternalServerError, &errorResponse{
		Code:   "INTERNAL_SERVER_ERROR",
		Detail: err.Error(),
	})
}
