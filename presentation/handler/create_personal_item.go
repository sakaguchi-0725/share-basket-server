package handler

import (
	"context"
	"net/http"
	"sharebasket/core"
	"sharebasket/usecase"

	"github.com/labstack/echo/v4"
)

type createPersonalItemRequest struct {
	Name       string `json:"name"`
	Status     string `json:"status"`
	CategoryID int64  `json:"categoryId"`
}

func NewCreatePersonalItem(usecase usecase.CreatePersonalItem, logger core.Logger) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req createPersonalItemRequest

		if err := c.Bind(&req); err != nil {
			logger.WithError(err).
				With("endpoint", c.Path()).
				With("method", c.Request().Method).
				Info("invalid request format")
			return core.NewInvalidError(err)
		}

		ctx := c.Request().Context()
		input, err := req.makeInput(ctx, logger)
		if err != nil {
			return err
		}

		err = usecase.Execute(ctx, input)
		if err != nil {
			return err
		}

		return c.NoContent(http.StatusNoContent)
	}
}

func (req *createPersonalItemRequest) makeInput(ctx context.Context, logger core.Logger) (usecase.CreatePersonalItemInput, error) {
	userID, err := core.GetUserID(ctx)
	if err != nil {
		logger.WithError(err).
			Info("failed to get user ID from context")
		return usecase.CreatePersonalItemInput{}, err
	}

	return usecase.CreatePersonalItemInput{
		UserID:     userID,
		Name:       req.Name,
		Status:     req.Status,
		CategoryID: req.CategoryID,
	}, nil
}
