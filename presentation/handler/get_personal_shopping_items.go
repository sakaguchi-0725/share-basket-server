package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"share-basket-server/core/apperr"
	contextKey "share-basket-server/core/context"
	"share-basket-server/presentation/response"
	"share-basket-server/usecase"
)

type (
	getPersonalShoppingItemsRequest struct {
		Status string `json:"status"`
	}

	getPersonalShoppingItemResponse struct {
		ID         uint   `json:"id"`
		Name       string `json:"name"`
		Status     string `json:"status"`
		CategoryID uint   `json:"categoryId"`
	}

	getPersonalShoppingItemsPresenter struct {
		w http.ResponseWriter
	}
)

func MakeGetPersonalShoppingItemsHandler(usecase usecase.GetPersonalShoppingItemsInputPort) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req getPersonalShoppingItemsRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Error(w, apperr.NewInvalidError(err))
			return
		}

		ctx := r.Context()
		err := usecase.Execute(ctx, req.makeInput(ctx), newGetPersonalShoppingItemsPresenter(w))
		if err != nil {
			response.Error(w, err)
		}
	}
}

func (req getPersonalShoppingItemsRequest) makeInput(ctx context.Context) usecase.GetPersonalShoppingItemsInput {
	return usecase.GetPersonalShoppingItemsInput{
		UserID: ctx.Value(contextKey.UserID).(string),
		Status: req.Status,
	}
}

func (presenter *getPersonalShoppingItemsPresenter) Render(
	ctx context.Context,
	output []usecase.GetPersonalShoppingItemOutput,
) error {
	res := make([]getPersonalShoppingItemResponse, len(output))

	for i, v := range output {
		res[i] = getPersonalShoppingItemResponse{
			ID:         v.ID,
			Name:       v.Name,
			Status:     v.Status,
			CategoryID: v.CategoryID,
		}
	}

	return response.StatusOK(presenter.w, res)
}

func newGetPersonalShoppingItemsPresenter(w http.ResponseWriter) usecase.GetPersonalShoppingItemsOutputPort {
	return &getPersonalShoppingItemsPresenter{w}
}
