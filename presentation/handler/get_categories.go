package handler

import (
	"net/http"
	"sharebasket/presentation/response"
	"sharebasket/usecase"
)

type getCategoriesResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func NewGetCategories(usecase usecase.GetCategories) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		out, err := usecase.Execute(r.Context())
		if err != nil {
			response.Error(w, err)
			return
		}

		response.StatusOK(w, makeGetCategoriesResponse(out))
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
