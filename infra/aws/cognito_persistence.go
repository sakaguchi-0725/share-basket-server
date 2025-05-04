package aws

import (
	"context"
	"errors"
	"fmt"
	"share-basket-server/core/apperr"
	"share-basket-server/core/util"
	"share-basket-server/domain"
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

func (c *cognito) Login(ctx context.Context, email string, password string) (string, error) {
	result, err := c.client.InitiateAuth(ctx, email, password)

	if err != nil {
		if errors.As(err, &notAuthorizedException) {
			return "", apperr.ErrUnauthenticated
		}

		if errors.As(err, &invalidParameterException) || errors.As(err, &invalidPasswordException) {
			return "", apperr.ErrInvalidData
		}

		return "", err
	}

	accessToken := result.AuthenticationResult.AccessToken
	if accessToken == nil {
		return "", apperr.ErrUnauthenticated
	}

	return util.Derefer(accessToken), nil
}

func (c *cognito) SignUp(ctx context.Context, email string, password string) (string, error) {
	result, err := c.client.SignUp(ctx, email, password)

	if err != nil {
		if errors.As(err, &usernameExistsException) {
			return "", apperr.ErrDuplicatedKey
		}
		if errors.As(err, &invalidParameterException) || errors.As(err, &invalidPasswordException) {
			return "", apperr.ErrInvalidData
		}
		return "", fmt.Errorf("failed to sign up: %w", err)
	}

	cognitoUID := result.UserSub
	if cognitoUID == nil {
		return "", apperr.ErrInvalidData
	}

	return util.Derefer(cognitoUID), nil
}

func (c *cognito) SignUpConfirm(ctx context.Context, email string, confirmationCode string) error {
	_, err := c.client.ConfirmSignUp(ctx, email, confirmationCode)

	if err != nil {
		if errors.As(err, &invalidParameterException) || errors.As(err, &codeMismatchException) {
			return apperr.ErrInvalidData
		}

		if errors.As(err, &expiredCodeException) {
			return apperr.ErrExpiredCodeException
		}

		return err
	}

	return nil
}

func (c *cognito) VerifyToken(ctx context.Context, token string) (string, error) {
	parsedToken, err := c.client.ParseToken(ctx, token)
	if err != nil {
		return "", apperr.ErrInvalidToken
	}

	if !parsedToken.Valid {
		return "", apperr.ErrInvalidToken
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return "", apperr.ErrInvalidToken
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		return "", apperr.ErrInvalidToken
	}

	if int64(exp) < time.Now().Unix() {
		return "", apperr.ErrTokenExpired
	}

	email, ok := claims["username"].(string)
	if !ok {
		return "", apperr.ErrInvalidToken
	}

	return email, nil
}

func (c *cognito) DeleteUser(ctx context.Context, email string) error {
	return c.client.AdminDeleteUser(ctx, email)
}

func NewCognitoPersistence(client CognitoClient) domain.Authenticator {
	return &cognito{client}
}
