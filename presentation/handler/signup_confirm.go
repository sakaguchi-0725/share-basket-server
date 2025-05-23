package handler

import (
	"net/http"
	"sharebasket/core"
	"sharebasket/usecase"

	"github.com/labstack/echo/v4"
)

type signUpConfirmRequest struct {
	Email            string `json:"email"`
	ConfirmationCode string `json:"confirmationCode"`
}

func NewSignUpConfirm(usecase usecase.SignUpConfirm) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req signUpConfirmRequest

		if err := c.Bind(&req); err != nil {
			return core.NewInvalidError(err)
		}

		err := usecase.Execute(c.Request().Context(), req.makeInput())
		if err != nil {
			return err
		}

		return c.NoContent(http.StatusNoContent)
	}
}

func (s *signUpConfirmRequest) makeInput() usecase.SignUpConfirmInput {
	return usecase.SignUpConfirmInput{
		Email:            s.Email,
		ConfirmationCode: s.ConfirmationCode,
	}
}
