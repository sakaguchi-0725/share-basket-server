package handler

import (
	"net/http"
	"sharebasket/core"
	"sharebasket/usecase"

	"github.com/labstack/echo/v4"
)

type getPersonalItemsResponse struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	Status     string `json:"status"`
	CategoryID int64  `json:"categoryId"`
}

func NewGetPersonalItems(usecase usecase.GetPersonalItems, logger core.Logger) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		userID, err := core.GetUserID(ctx)
		if err != nil {
			logger.WithError(err).
				Info("failed to get user ID from context")
			return err
		}

		status := c.QueryParam("status")

		out, err := usecase.Execute(ctx, makeGetPersonalItemsInput(userID, status))
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, makeGetPersonalItemsResponse(out))
	}
}

func makeGetPersonalItemsInput(userID, status string) usecase.GetPersonalItemsInput {
	return usecase.GetPersonalItemsInput{
		UserID: userID,
		Status: status,
	}
}

func makeGetPersonalItemsResponse(out []usecase.GetPersonalItemsOutput) []getPersonalItemsResponse {
	res := make([]getPersonalItemsResponse, len(out))

	for i, v := range out {
		res[i] = getPersonalItemsResponse{
			ID:         v.ID,
			Name:       v.Name,
			Status:     v.Status,
			CategoryID: v.CategoryID,
		}
	}

	return res
}
