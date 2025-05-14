package handler

import (
	"encoding/json"
	"net/http"
	"sharebasket/core"
	"sharebasket/presentation/response"
	"sharebasket/usecase"
)

type signUpConfirmRequest struct {
	Email            string `json:"email"`
	ConfirmationCode string `json:"confirmationCode"`
}

func NewSignUpConfirm(usecase usecase.SignUpConfirm) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req signUpConfirmRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Error(w, core.NewInvalidError(err))
			return
		}

		err := usecase.Execute(r.Context(), req.makeInput())
		if err != nil {
			response.Error(w, err)
			return
		}

		response.NoContent(w)
	}
}

func (s *signUpConfirmRequest) makeInput() usecase.SignUpConfirmInput {
	return usecase.SignUpConfirmInput{
		Email:            s.Email,
		ConfirmationCode: s.ConfirmationCode,
	}
}
