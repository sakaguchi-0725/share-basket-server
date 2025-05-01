package persistence

import (
	"context"
	"crypto/hmac"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"share-basket-server/core/apperr"
	appConfig "share-basket-server/core/config"
	"share-basket-server/core/util"
	"share-basket-server/personal/domain/repository"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"

	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
)

type cognito struct {
	client       *cognitoidentityprovider.Client
	clientID     string
	clientSecret string
	userPoolID   string
	jwksURL      string
	keys         map[string]*rsa.PublicKey
	mu           sync.Mutex
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
	secretHash := c.genSecretHash(email)
	result, err := c.client.InitiateAuth(ctx, &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: types.AuthFlowTypeUserPasswordAuth,
		ClientId: aws.String(c.clientID),
		AuthParameters: map[string]string{
			"USERNAME":    email,
			"PASSWORD":    password,
			"SECRET_HASH": secretHash,
		},
	})

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
	secretHash := c.genSecretHash(email)

	result, err := c.client.SignUp(ctx, &cognitoidentityprovider.SignUpInput{
		ClientId:   aws.String(c.clientID),
		Username:   aws.String(email),
		Password:   aws.String(password),
		SecretHash: aws.String(secretHash),
	})

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
	secretHash := c.genSecretHash(email)

	_, err := c.client.ConfirmSignUp(ctx, &cognitoidentityprovider.ConfirmSignUpInput{
		ClientId:         aws.String(c.clientID),
		Username:         aws.String(email),
		ConfirmationCode: aws.String(confirmationCode),
		SecretHash:       aws.String(secretHash),
	})

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

// VerifyToken implements repository.Authenticator.
func (c *cognito) VerifyToken(ctx context.Context, token string) (string, error) {
	panic("unimplemented")
}

func (c *cognito) DeleteUser(ctx context.Context, email string) error {
	_, err := c.client.AdminDeleteUser(ctx, &cognitoidentityprovider.AdminDeleteUserInput{
		UserPoolId: aws.String(c.userPoolID),
		Username:   aws.String(email),
	})

	return err
}

// Cognito認証に必要なSecretHashを生成する
func (c *cognito) genSecretHash(username string) string {
	mac := hmac.New(sha256.New, []byte(c.clientSecret))
	mac.Write([]byte(username + c.clientID))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func NewCognito(ctx context.Context, cfg appConfig.AWS) (repository.Authenticator, error) {
	awsCfg, err := awsConfig.LoadDefaultConfig(ctx, awsConfig.WithRegion(cfg.Region))
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	return &cognito{
		client:       cognitoidentityprovider.NewFromConfig(awsCfg),
		clientID:     cfg.CognitoClientID,
		clientSecret: cfg.CognitoClientSecret,
		userPoolID:   cfg.CognitoUserPoolID,
		jwksURL: fmt.Sprintf(
			"https://cognito-idp.%s.amazonaws.com/%s/.well-known/jwks.json",
			cfg.Region, cfg.CognitoUserPoolID,
		),
	}, nil
}
