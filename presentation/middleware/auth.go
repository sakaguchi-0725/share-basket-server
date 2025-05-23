package middleware

import (
	"context"
	"net/http"
	"sharebasket/core"
	"sharebasket/usecase"

	"github.com/labstack/echo/v4"
)

func Auth(v usecase.VerifyToken) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := c.Request().Context()

			cookie, err := c.Cookie("access_token")
			if err != nil {
				if err == http.ErrNoCookie {
					return core.NewAppError(core.ErrUnauthorized, err)
				}

				return err
			}

			userID, err := v.Execute(ctx, cookie.Value)
			if err != nil {
				return err
			}

			newCtx := context.WithValue(ctx, core.UserIDKey, userID)
			c.SetRequest(c.Request().WithContext(newCtx))
			return next(c)
		}
	}
}
