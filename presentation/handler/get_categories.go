package handler

import (
	"net/http"
	"sharebasket/usecase"

	"github.com/labstack/echo/v4"
)

type getCategoriesResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func NewGetCategories(usecase usecase.GetCategories) echo.HandlerFunc {
	return func(c echo.Context) error {
		out, err := usecase.Execute(c.Request().Context())
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, makeGetCategoriesResponse(out))
	}
}

func makeGetCategoriesResponse(out []usecase.GetCategoriesOutput) []getCategoriesResponse {
	res := make([]getCategoriesResponse, len(out))

	for i, v := range out {
		res[i] = getCategoriesResponse{
			ID:   v.ID,
			Name: v.Name,
		}
	}

	return res
}
