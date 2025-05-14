package handler

import (
	"encoding/json"
	"net/http"
	"sharebasket/core"
	"sharebasket/presentation/response"
	"sharebasket/usecase"
)

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewLogin(usecase usecase.Login, logger core.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req loginRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			logger.WithError(err).
				With("endpoint", r.URL.Path).
				With("method", r.Method).
				Info("invalid request format")
			response.Error(w, core.NewInvalidError(err))
			return
		}

		output, err := usecase.Execute(r.Context(), req.makeInput())
		if err != nil {
			response.Error(w, err)
			return
		}

		setCookies(w, r, output)
		response.NoContent(w)
	}
}

func (l *loginRequest) makeInput() usecase.LoginInput {
	return usecase.LoginInput{
		Email:    l.Email,
		Password: l.Password,
	}
}

func setCookies(w http.ResponseWriter, r *http.Request, o usecase.LoginOutput) {
	accessToken := &http.Cookie{
		Name:     "access_token",
		Value:    o.AccessToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   r.TLS != nil,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   86400,
	}

	refreshToken := &http.Cookie{
		Name:     "refresh_token",
		Value:    o.RefreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   r.TLS != nil,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   86400,
	}

	http.SetCookie(w, accessToken)
	http.SetCookie(w, refreshToken)
}
