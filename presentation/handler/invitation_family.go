package handler

import (
	"net/http"
	"sharebasket/core"
	"sharebasket/presentation/response"
	"sharebasket/usecase"
)

type invitationFamilyResponse struct {
	Token string `json:"token"`
}

func NewInvitationFamily(usecase usecase.InvitationFamily, logger core.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		userID, err := core.GetUserID(ctx)
		if err != nil {
			logger.WithError(err).
				Info("failed to get user ID from context")
			response.Error(w, err)
			return
		}

		token, err := usecase.Execute(ctx, userID)
		if err != nil {
			response.Error(w, err)
			return
		}

		response.StatusOK(w, makeInvitationFamilyResponse(token))
	}
}

func makeInvitationFamilyResponse(token string) invitationFamilyResponse {
	return invitationFamilyResponse{
		Token: token,
	}
}
