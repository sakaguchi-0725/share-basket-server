package handler

import (
	"net/http"
	"sharebasket/core"
	"sharebasket/usecase"

	"github.com/labstack/echo/v4"
)

type getAccountResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func NewGetAccount(usecase usecase.GetAccount, logger core.Logger) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		userID, err := core.GetUserID(ctx)
		if err != nil {
			logger.WithError(err).
				Info("failed to get user ID from context")
			return err
		}

		out, err := usecase.Execute(ctx, userID)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, makeGetAccountResponse(out))
	}
}

func makeGetAccountResponse(out usecase.GetAccountOutput) getAccountResponse {
	return getAccountResponse{
		ID:   out.ID,
		Name: out.Name,
	}
}
