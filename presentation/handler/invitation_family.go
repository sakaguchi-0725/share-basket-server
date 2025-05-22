package handler

import (
	"net/http"
	"sharebasket/core"
	"sharebasket/usecase"

	"github.com/labstack/echo/v4"
)

type invitationFamilyResponse struct {
	Token string `json:"token"`
}

func NewInvitationFamily(usecase usecase.InvitationFamily, logger core.Logger) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		userID, err := core.GetUserID(ctx)
		if err != nil {
			logger.WithError(err).
				Info("failed to get user ID from context")
			return err
		}

		token, err := usecase.Execute(ctx, userID)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, makeInvitationFamilyResponse(token))
	}
}

func makeInvitationFamilyResponse(token string) invitationFamilyResponse {
	return invitationFamilyResponse{
		Token: token,
	}
}
