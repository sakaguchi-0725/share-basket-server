package handler

import (
	"errors"
	"net/http"
	"sharebasket/core"
	"sharebasket/usecase"

	"github.com/labstack/echo/v4"
)

type (
	getFamilyItemsResponse struct {
		ID         int64  `json:"id"`
		Name       string `json:"name"`
		Status     string `json:"status"`
		CategoryID int64  `json:"categoryId"`
	}
)

func NewGetFamilyItems(usecase usecase.GetFamilyItems) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		userID, err := core.GetUserID(ctx)
		if err != nil {
			return err
		}

		familyID := c.Param("familyID")
		if familyID == "" {
			return core.NewInvalidError(errors.New("family ID is required"))
		}

		status := c.QueryParam("status")

		out, err := usecase.Execute(ctx, makeGetFamilyItemsInput(userID, familyID, status))
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, makeGetFamilyItemsResponse(out))
	}
}

func makeGetFamilyItemsInput(userID, familyID, status string) usecase.GetFamilyItemsInput {
	return usecase.GetFamilyItemsInput{
		UserID:   userID,
		FamilyID: familyID,
		Status:   status,
	}
}

func makeGetFamilyItemsResponse(items []usecase.GetFamilyItemOutput) []getFamilyItemsResponse {
	res := make([]getFamilyItemsResponse, len(items))
	for i, item := range items {
		res[i] = getFamilyItemsResponse{
			ID:         item.ID,
			Name:       item.Name,
			Status:     item.Status,
			CategoryID: item.CategoryID,
		}
	}
	return res
}
