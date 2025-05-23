package handler

import (
	"net/http"
	"sharebasket/core"
	"sharebasket/usecase"

	"github.com/labstack/echo/v4"
)

type createFamilyRequest struct {
	Name string `json:"name"`
}

func NewCreateFamily(usecase usecase.CreateFamily) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req createFamilyRequest

		if err := c.Bind(&req); err != nil {
			return core.NewInvalidError(err)
		}

		input, err := req.makeInput(c)
		if err != nil {
			return err
		}

		if err := usecase.Execute(c.Request().Context(), input); err != nil {
			return err
		}

		return c.NoContent(http.StatusNoContent)
	}
}

func (req createFamilyRequest) makeInput(c echo.Context) (usecase.CreateFamilyInput, error) {
	id, err := core.GetUserID(c.Request().Context())
	if err != nil {
		return usecase.CreateFamilyInput{}, err
	}

	return usecase.CreateFamilyInput{
		Name:   req.Name,
		UserID: id,
	}, nil
}
