package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sharebasket/core"
	"sharebasket/presentation/response"
	"sharebasket/usecase"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type updatePersonalItemRequest struct {
	Name       string `json:"name"`
	Status     string `json:"status"`
	CategoryID int64  `json:"categoryId"`
}

func NewUpdatePersonalItem(usecase usecase.UpdatePersonalItem) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")

		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			response.Error(w, core.NewInvalidError(
				fmt.Errorf("invalid item id: %w", err),
			))
			return
		}

		var req updatePersonalItemRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Error(w, core.NewInvalidError(core.ErrInvalidData))
			return
		}

		ctx := r.Context()

		in, err := req.makeUpdatePersonalItemInput(ctx, id)
		if err != nil {
			response.Error(w, err)
			return
		}

		if err := usecase.Execute(ctx, in); err != nil {
			response.Error(w, err)
			return
		}

		response.NoContent(w)
	}
}

func (req updatePersonalItemRequest) makeUpdatePersonalItemInput(ctx context.Context, id int64) (usecase.UpdatePersonalItemInput, error) {
	userID, err := core.GetUserID(ctx)
	if err != nil {
		return usecase.UpdatePersonalItemInput{}, err
	}

	return usecase.UpdatePersonalItemInput{
		ID:         id,
		Name:       req.Name,
		Status:     req.Status,
		CategoryID: req.CategoryID,
		UserID:     userID,
	}, nil
}
