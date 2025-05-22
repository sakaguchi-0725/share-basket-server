package middleware

import (
	"net/http"
	"sharebasket/core"

	"github.com/labstack/echo/v4"
)

type errorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func Error() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)
			if err == nil {
				return nil
			}

			if appErr, ok := err.(*core.AppError); ok {
				var status int

				switch appErr.Code() {
				case core.ErrBadRequest, core.ErrExpiredCode:
					status = http.StatusBadRequest
				case core.ErrEmailAlreadyExists:
					status = http.StatusConflict
				case core.ErrForbidden:
					status = http.StatusForbidden
				case core.ErrUnauthorized, core.ErrExpiredToken:
					status = http.StatusUnauthorized
				case core.ErrNotFound:
					status = http.StatusNotFound
				default:
					status = http.StatusInternalServerError
				}

				return c.JSON(status, &errorResponse{
					Code:    appErr.Code().String(),
					Message: appErr.Error(),
				})
			}

			return c.JSON(http.StatusInternalServerError, &errorResponse{
				Code:    "INTERNAL_SERVER_ERROR",
				Message: "予期しないエラーが発生しました",
			})
		}
	}
}
