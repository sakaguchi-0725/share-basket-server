package handler

import (
	"fmt"
	"net/http"
	"sharebasket/core"
	"sharebasket/presentation/response"
	"sharebasket/usecase"
	"sharebasket/usecase/input"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func NewDeletePersonalItem(usecase usecase.PersonalItem, logger core.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")

		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			logger.WithError(err).
				With("item_id", idStr).
				With("endpoint", r.URL.Path).
				With("method", r.Method).
				Info("invalid item id")
			response.Error(w, core.NewInvalidError(
				fmt.Errorf("invalid item id: %w", err),
			))
			return
		}

		ctx := r.Context()

		userID, err := core.GetUserID(ctx)
		if err != nil {
			logger.WithError(err).
				Info("failed to get user ID from context")
			response.Error(w, core.NewInvalidError(err))
			return
		}

		err = usecase.Delete(ctx, makeDeletePersonalItemInput(id, userID))
		if err != nil {
			response.Error(w, err)
			return
		}

		response.NoContent(w)
	}
}

func makeDeletePersonalItemInput(id int64, userID string) input.DeletePersonalItem {
	return input.DeletePersonalItem{
		ID:     id,
		UserID: userID,
	}
}
