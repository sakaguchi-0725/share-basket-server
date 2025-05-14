package handler

import (
	"encoding/json"
	"net/http"
	"sharebasket/core"
	"sharebasket/presentation/response"
	"sharebasket/usecase"
)

type signUpRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewSignUp(usecase usecase.SignUp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req signUpRequest

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

func (s *signUpRequest) makeInput() usecase.SignUpInput {
	return usecase.SignUpInput{
		Name:     s.Name,
		Email:    s.Email,
		Password: s.Password,
	}
}
