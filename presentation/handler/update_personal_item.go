package handler

import (
	"context"
	"fmt"
	"net/http"
	"sharebasket/core"
	"sharebasket/usecase"
	"strconv"

	"github.com/labstack/echo/v4"
)

type updatePersonalItemRequest struct {
	Name       string `json:"name"`
	Status     string `json:"status"`
	CategoryID int64  `json:"categoryId"`
}

func NewUpdatePersonalItem(usecase usecase.UpdatePersonalItem) echo.HandlerFunc {
	return func(c echo.Context) error {
		idStr := c.Param("id")

		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return core.NewInvalidError(
				fmt.Errorf("invalid item id: %w", err),
			)
		}

		var req updatePersonalItemRequest
		if err := c.Bind(&req); err != nil {
			return core.NewInvalidError(core.ErrInvalidData)
		}

		ctx := c.Request().Context()

		in, err := req.makeUpdatePersonalItemInput(ctx, id)
		if err != nil {
			return err
		}

		if err := usecase.Execute(ctx, in); err != nil {
			return err
		}

		return c.NoContent(http.StatusNoContent)
	}
}

func (req updatePersonalItemRequest) makeUpdatePersonalItemInput(ctx context.Context, id int64) (usecase.UpdatePersonalItemInput, error) {
	userID, err := core.GetUserID(ctx)
	if err != nil {
		return usecase.UpdatePersonalItemInput{}, err
	}

	return usecase.UpdatePersonalItemInput{
		ID:         id,
		Name:       req.Name,
		Status:     req.Status,
		CategoryID: req.CategoryID,
		UserID:     userID,
	}, nil
}
