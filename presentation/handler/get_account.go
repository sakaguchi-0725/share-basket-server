package handler

import (
	"net/http"
	"sharebasket/core"
	"sharebasket/presentation/response"
	"sharebasket/usecase"
)

type getAccountResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func NewGetAccount(usecase usecase.GetAccount, logger core.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userID, err := core.GetUserID(ctx)
		if err != nil {
			logger.WithError(err).
				Info("failed to get user ID from context")
			response.Error(w, err)
			return
		}

		out, err := usecase.Execute(ctx, userID)
		if err != nil {
			response.Error(w, err)
			return
		}

		response.StatusOK(w, makeGetAccountResponse(out))
	}
}

func makeGetAccountResponse(out usecase.GetAccountOutput) getAccountResponse {
	return getAccountResponse{
		ID:   out.ID,
		Name: out.Name,
	}
}
