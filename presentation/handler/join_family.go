package handler

import (
	"errors"
	"net/http"
	"sharebasket/core"
	"sharebasket/usecase"

	"github.com/labstack/echo/v4"
)

func NewJoinFamily(uc usecase.JoinFamily) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Param("token")
		if token == "" {
			return core.NewInvalidError(errors.New("token is empty"))
		}

		ctx := c.Request().Context()

		userID, err := core.GetUserID(ctx)
		if err != nil {
			return err
		}

		err = uc.Execute(ctx, usecase.JoinFamilyInput{
			UserID: userID,
			Token:  token,
		})

		if err != nil {
			return err
		}

		return c.NoContent(http.StatusNoContent)
	}
}
