package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func NewHealthCheck() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"message": "ok",
		})
	}
}
