package auth

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sharebasket/core"
	"sharebasket/domain/repository"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/golang-jwt/jwt/v5"
)

type cognito struct {
	client CognitoClient
}

var (
	usernameExistsException   = &types.UsernameExistsException{}
	invalidPasswordException  = &types.InvalidPasswordException{}
	invalidParameterException = &types.InvalidParameterException{}
	expiredCodeException      = &types.ExpiredCodeException{}
	codeMismatchException     = &types.CodeMismatchException{}
	notAuthorizedException    = &types.NotAuthorizedException{}
)

func (c *cognito) Login(ctx context.Context, email string, password string) (accessToken, refreshToken string, err error) {
	result, err := c.client.InitiateAuth(ctx, email, password)

	if err != nil {
		if errors.As(err, &notAuthorizedException) {
			return "", "", core.NewAppError(core.ErrUnauthorized, ErrLoginFailed).WithMessage("ログインに失敗しました")
		}

		if errors.As(err, &invalidParameterException) || errors.As(err, &invalidPasswordException) {
			return "", "", core.NewInvalidError(ErrInvalidInput)
		}

		return "", "", err
	}

	accessTokenPtr := result.AuthenticationResult.AccessToken
	refreshTokenPtr := result.AuthenticationResult.RefreshToken
	if accessTokenPtr == nil || refreshTokenPtr == nil {
		return "", "", core.NewAppError(core.ErrUnauthorized, ErrLoginFailed)
	}

	return core.Derefer(accessTokenPtr), core.Derefer(refreshTokenPtr), nil
}

func (c *cognito) SignUp(ctx context.Context, email string, password string) (string, error) {
	result, err := c.client.SignUp(ctx, email, password)
	log.Println(err)

	if err != nil {
		if errors.As(err, &usernameExistsException) {
			return "", core.NewInvalidError(ErrEmailExists).WithMessage("このメールアドレスは使用できません")
		}
		if errors.As(err, &invalidParameterException) || errors.As(err, &invalidPasswordException) {
			return "", core.NewInvalidError(ErrInvalidInput)
		}
		return "", fmt.Errorf("failed to sign up: %w", err)
	}

	cognitoUID := result.UserSub
	if cognitoUID == nil {
		return "", core.NewInvalidError(ErrInvalidInput)
	}

	return core.Derefer(cognitoUID), nil
}

func (c *cognito) SignUpConfirm(ctx context.Context, email string, confirmationCode string) error {
	_, err := c.client.ConfirmSignUp(ctx, email, confirmationCode)

	if err != nil {
		if errors.As(err, &invalidParameterException) || errors.As(err, &codeMismatchException) {
			return core.NewInvalidError(ErrInvalidInput)
		}

		if errors.As(err, &expiredCodeException) {
			return core.NewInvalidError(ErrExpiredCode).WithMessage("確認コードの有効期限が切れてます")
		}

		return err
	}

	return nil
}

func (c *cognito) VerifyToken(ctx context.Context, token string) (string, error) {
	parsedToken, err := c.client.ParseToken(ctx, token)
	if err != nil {
		return "", core.NewUnauthorizedError(ErrInvalidAccessToken)
	}

	if !parsedToken.Valid {
		return "", core.NewUnauthorizedError(ErrInvalidAccessToken)
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return "", core.NewUnauthorizedError(ErrInvalidAccessToken)
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		return "", core.NewUnauthorizedError(ErrInvalidAccessToken)
	}

	if int64(exp) < time.Now().Unix() {
		return "", core.NewUnauthorizedError(ErrExpiredCode)
	}

	email, ok := claims["username"].(string)
	if !ok {
		return "", core.NewUnauthorizedError(ErrExpiredCode)
	}

	return email, nil
}

func (c *cognito) DeleteUser(ctx context.Context, email string) error {
	return c.client.AdminDeleteUser(ctx, email)
}

func NewCognito(client CognitoClient) repository.Authenticator {
	return &cognito{client}
}
