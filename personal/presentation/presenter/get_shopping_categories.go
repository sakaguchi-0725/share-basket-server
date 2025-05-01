package presenter

import (
	"context"
	"net/http"
	"share-basket-server/personal/presentation/response"
	"share-basket-server/personal/usecase/output"
)

type getShoppingCategories struct {
	w http.ResponseWriter
}

type getShoppingCategoriesResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

func (g *getShoppingCategories) Render(ctx context.Context, out output.GetShoppingCategories) error {
	return response.JSON(g.w, http.StatusOK, g.makeResponse(out))
}

func (g *getShoppingCategories) makeResponse(outputs output.GetShoppingCategories) []getShoppingCategoriesResponse {
	res := make([]getShoppingCategoriesResponse, len(outputs))

	for i, v := range outputs {
		res[i] = getShoppingCategoriesResponse{
			ID:   v.ID,
			Name: v.Name,
		}
	}

	return res
}

func NewGetShoppingCategories(w http.ResponseWriter) output.GetShoppingCategoriesPort {
	return &getShoppingCategories{w}
}
