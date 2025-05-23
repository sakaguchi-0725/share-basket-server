package handler

import (
	"fmt"
	"net/http"
	"sharebasket/core"
	"sharebasket/usecase"
	"strconv"

	"github.com/labstack/echo/v4"
)

func NewDeletePersonalItem(usecase usecase.DeletePersonalItem) echo.HandlerFunc {
	return func(c echo.Context) error {
		idStr := c.Param("id")

		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return core.NewInvalidError(
				fmt.Errorf("invalid item id: %w", err),
			)
		}

		ctx := c.Request().Context()

		userID, err := core.GetUserID(ctx)
		if err != nil {
			return core.NewInvalidError(err)
		}

		err = usecase.Execute(ctx, makeDeletePersonalItemInput(id, userID))
		if err != nil {
			return err
		}

		return c.NoContent(http.StatusNoContent)
	}
}

func makeDeletePersonalItemInput(id int64, userID string) usecase.DeletePersonalItemInput {
	return usecase.DeletePersonalItemInput{
		ID:     id,
		UserID: userID,
	}
}
