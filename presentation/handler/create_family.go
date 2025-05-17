package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"sharebasket/core"
	"sharebasket/presentation/response"
	"sharebasket/usecase"
)

type createFamilyRequest struct {
	Name string `json:"name"`
}

func NewCreateFamily(usecase usecase.CreateFamily, logger core.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createFamilyRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			logger.WithError(err).
				With("endpoint", r.URL.Path).
				With("method", r.Method).
				Info("invalid request format")
			response.Error(w, core.NewInvalidError(err))
			return
		}

		ctx := r.Context()

		input, err := req.makeInput(ctx)
		if err != nil {
			logger.WithError(err).
				Info("failed to get user ID from context")

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

func (req createFamilyRequest) makeInput(ctx context.Context) (usecase.CreateFamilyInput, error) {
	id, err := core.GetUserID(ctx)
	if err != nil {
		return usecase.CreateFamilyInput{}, err
	}

	return usecase.CreateFamilyInput{
		Name:   req.Name,
		UserID: id,
	}, nil
}
