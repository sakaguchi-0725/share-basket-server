package middleware

import (
	"fmt"
	"net/http"
	"sharebasket/core"

	"github.com/labstack/echo/v4"
)

type errorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func Error(logger core.Logger) echo.MiddlewareFunc {
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

				errLogger := logger.
					With("error_code", appErr.Code().String()).
					With("stack_trace", appErr.StackTrace()).
					With("status", status).
					With("path", c.Request().URL.Path)

				logMsg := fmt.Sprintf("AppError: %s", appErr.Error())

				if status >= http.StatusInternalServerError {
					errLogger.Error(logMsg)
				} else {
					errLogger.Warn(logMsg)
				}

				return c.JSON(status, &errorResponse{
					Code:    appErr.Code().String(),
					Message: appErr.Message(),
				})
			}

			logger.WithError(err).
				With("path", c.Request().URL.Path).
				Error(fmt.Sprintf("予期しないエラー: %v", err))

			return c.JSON(http.StatusInternalServerError, &errorResponse{
				Code:    "INTERNAL_SERVER_ERROR",
				Message: "予期しないエラーが発生しました",
			})
		}
	}
}
