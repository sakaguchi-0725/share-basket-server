package handler

import (
	"net/http"
	"sharebasket/core"
	"sharebasket/presentation/response"
	"sharebasket/usecase"
)

type getPersonalItemsResponse struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	Status     string `json:"status"`
	CategoryID int64  `json:"categoryId"`
}

func NewGetPersonalItems(usecase usecase.GetPersonalItems) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		userID, err := core.GetUserID(ctx)
		if err != nil {
			response.Error(w, err)
			return
		}

		status := r.URL.Query().Get("status")

		out, err := usecase.Execute(ctx, makeGetPersonalItemsInput(userID, status))
		if err != nil {
			response.Error(w, err)
			return
		}

		response.StatusOK(w, makeGetPersonalItemsResponse(out))
	}
}

func makeGetPersonalItemsInput(userID, status string) usecase.GetPersonalItemsInput {
	return usecase.GetPersonalItemsInput{
		UserID: userID,
		Status: status,
	}
}

func makeGetPersonalItemsResponse(out []usecase.GetPersonalItemsOutput) []getPersonalItemsResponse {
	res := make([]getPersonalItemsResponse, len(out))

	for i, v := range out {
		res[i] = getPersonalItemsResponse{
			ID:         v.ID,
			Name:       v.Name,
			Status:     v.Status,
			CategoryID: v.CategoryID,
		}
	}

	return res
}
