package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"sharebasket/core"
	"sharebasket/presentation/response"
	"sharebasket/usecase"
)

type createPersonalItemRequest struct {
	Name       string `json:"name"`
	Status     string `json:"status"`
	CategoryID int64  `json:"categoryId"`
}

func NewCreatePersonalItem(usecase usecase.CreatePersonalItem) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createPersonalItemRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Error(w, core.NewInvalidError(err))
			return
		}

		ctx := r.Context()
		input, err := req.makeInput(ctx)
		if err != nil {
			response.Error(w, err)
			return
		}

		err = usecase.Execute(ctx, input)
		if err != nil {
			response.Error(w, err)
			return
		}

		response.NoContent(w)
	}
}

func (req *createPersonalItemRequest) makeInput(ctx context.Context) (usecase.CreatePersonalItemInput, error) {
	userID, err := core.GetUserID(ctx)
	if err != nil {
		return usecase.CreatePersonalItemInput{}, err
	}

	return usecase.CreatePersonalItemInput{
		UserID:     userID,
		Name:       req.Name,
		Status:     req.Status,
		CategoryID: req.CategoryID,
	}, nil
}
