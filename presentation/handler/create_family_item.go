package handler

import (
	"net/http"
	"sharebasket/core"
	"sharebasket/usecase"

	"github.com/labstack/echo/v4"
)

type createFamilyItemRequest struct {
	Name       string `json:"name"`
	Status     string `json:"status"`
	CategoryID int64  `json:"categoryId"`
}

func NewCreateFamilyItem(uc usecase.CreateFamilyItem) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req createFamilyItemRequest

		if err := c.Bind(&req); err != nil {
			return core.NewInvalidError(err)
		}

		ctx := c.Request().Context()

		userID, err := core.GetUserID(ctx)
		if err != nil {
			return err
		}

		err = uc.Execute(ctx, req.makeInput(userID))
		if err != nil {
			return err
		}

		return c.NoContent(http.StatusNoContent)
	}
}

func (req createFamilyItemRequest) makeInput(userID string) usecase.CreateFamilyItemInput {
	return usecase.CreateFamilyItemInput{
		Name:       req.Name,
		Status:     req.Status,
		CategoryID: req.CategoryID,
		UserID:     userID,
	}
}
