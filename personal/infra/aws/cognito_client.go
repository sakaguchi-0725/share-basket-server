//go:generate mockgen -destination=../../mock/$GOPACKAGE/$GOFILE . CognitoClient
package aws

import (
	"context"
	"crypto/hmac"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	appConfig "share-basket-server/core/config"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
)

type (
	CognitoClient interface {
		InitiateAuth(ctx context.Context, email, password string) (*cognitoidentityprovider.InitiateAuthOutput, error)
		SignUp(ctx context.Context, email, password string) (*cognitoidentityprovider.SignUpOutput, error)
		ConfirmSignUp(ctx context.Context, email, confirmationCode string) (*cognitoidentityprovider.ConfirmSignUpOutput, error)
		AdminDeleteUser(ctx context.Context, email string) error
	}

	cognitoClient struct {
		sdk        *cognitoidentityprovider.Client
		id         string
		secret     string
		userPoolID string
		jwksURL    string
		keys       map[string]*rsa.PublicKey
		mu         sync.Mutex
	}
)

func (client *cognitoClient) AdminDeleteUser(ctx context.Context, email string) error {
	_, err := client.sdk.AdminDeleteUser(ctx, &cognitoidentityprovider.AdminDeleteUserInput{
		UserPoolId: aws.String(client.userPoolID),
		Username:   aws.String(email),
	})

	return err
}

func (client *cognitoClient) ConfirmSignUp(ctx context.Context, email, confirmationCode string) (*cognitoidentityprovider.ConfirmSignUpOutput, error) {
	return client.sdk.ConfirmSignUp(ctx, &cognitoidentityprovider.ConfirmSignUpInput{
		ClientId:         aws.String(client.id),
		Username:         aws.String(email),
		ConfirmationCode: aws.String(confirmationCode),
		SecretHash:       aws.String(client.genSecretHash(email)),
	})
}

func (client *cognitoClient) InitiateAuth(ctx context.Context, email, password string) (*cognitoidentityprovider.InitiateAuthOutput, error) {
	return client.sdk.InitiateAuth(ctx, &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: types.AuthFlowTypeUserPasswordAuth,
		ClientId: aws.String(client.id),
		AuthParameters: map[string]string{
			"USERNAME":    email,
			"PASSWORD":    password,
			"SECRET_HASH": client.genSecretHash(email),
		},
	})
}

func (client *cognitoClient) SignUp(ctx context.Context, email, password string) (*cognitoidentityprovider.SignUpOutput, error) {
	return client.sdk.SignUp(ctx, &cognitoidentityprovider.SignUpInput{
		ClientId:   aws.String(client.id),
		Username:   aws.String(email),
		Password:   aws.String(password),
		SecretHash: aws.String(client.genSecretHash(email)),
	})
}

// Cognito認証に必要なSecretHashを生成する
func (client *cognitoClient) genSecretHash(username string) string {
	mac := hmac.New(sha256.New, []byte(client.secret))
	mac.Write([]byte(username + client.id))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func NewCognitoClient(ctx context.Context, cfg appConfig.AWS) (CognitoClient, error) {
	awsCfg, err := awsConfig.LoadDefaultConfig(ctx, awsConfig.WithRegion(cfg.Region))
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	return &cognitoClient{
		sdk:        cognitoidentityprovider.NewFromConfig(awsCfg),
		id:         cfg.CognitoClientID,
		secret:     cfg.CognitoClientSecret,
		userPoolID: cfg.CognitoUserPoolID,
		jwksURL: fmt.Sprintf(
			"https://cognito-idp.%s.amazonaws.com/%s/.well-known/jwks.json",
			cfg.Region, cfg.CognitoUserPoolID,
		),
	}, nil
}
