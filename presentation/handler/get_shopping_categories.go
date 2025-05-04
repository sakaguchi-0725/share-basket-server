package handler

import (
	"context"
	"net/http"
	"share-basket-server/presentation/response"
	"share-basket-server/usecase"
)

type (
	getShoppingCategoriesResponse struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	}

	getShoppingCategoriesPresenter struct {
		w http.ResponseWriter
	}
)

func MakeGetShoppingCategoriesHandler(usecase usecase.GetShoppingCategoriesInputPort) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := usecase.Execute(r.Context(), NewGetShoppingCategoriesPresenter(w))
		if err != nil {
			response.Error(w, err)
		}
	}
}

func (presenter *getShoppingCategoriesPresenter) Render(ctx context.Context, output []usecase.GetShoppingCategoryOutput) error {
	res := make([]getShoppingCategoriesResponse, len(output))

	for i, v := range output {
		res[i] = getShoppingCategoriesResponse{
			ID:   v.ID,
			Name: v.Name,
		}
	}

	return response.StatusOK(presenter.w, res)
}

func NewGetShoppingCategoriesPresenter(w http.ResponseWriter) usecase.GetShoppingCategoriesOutputPort {
	return &getShoppingCategoriesPresenter{w}
}
