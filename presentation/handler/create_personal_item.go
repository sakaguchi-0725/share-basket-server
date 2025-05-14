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

func NewCreatePersonalItem(usecase usecase.CreatePersonalItem, logger core.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createPersonalItemRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			logger.WithError(err).
				With("endpoint", r.URL.Path).
				With("method", r.Method).
				Info("invalid request format")
			response.Error(w, core.NewInvalidError(err))
			return
		}

		ctx := r.Context()
		input, err := req.makeInput(ctx, logger)
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

func (req *createPersonalItemRequest) makeInput(ctx context.Context, logger core.Logger) (usecase.CreatePersonalItemInput, error) {
	userID, err := core.GetUserID(ctx)
	if err != nil {
		logger.WithError(err).
			Info("failed to get user ID from context")
		return usecase.CreatePersonalItemInput{}, err
	}

	return usecase.CreatePersonalItemInput{
		UserID:     userID,
		Name:       req.Name,
		Status:     req.Status,
		CategoryID: req.CategoryID,
	}, nil
}
