package handler

import (
	"net/http"
	"sharebasket/core"
	"sharebasket/usecase"

	"github.com/labstack/echo/v4"
)

type signUpRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewSignUp(usecase usecase.SignUp, logger core.Logger) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req signUpRequest

		if err := c.Bind(&req); err != nil {
			logger.WithError(err).
				With("endpoint", c.Path()).
				With("method", c.Request().Method).
				Info("invalid request format")
			return core.NewInvalidError(err)
		}

		err := usecase.Execute(c.Request().Context(), req.makeInput())
		if err != nil {
			return err
		}

		return c.NoContent(http.StatusNoContent)
	}
}

func (s *signUpRequest) makeInput() usecase.SignUpInput {
	return usecase.SignUpInput{
		Name:     s.Name,
		Email:    s.Email,
		Password: s.Password,
	}
}
