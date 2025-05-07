//go:generate mockgen -destination=../../test/mock/$GOPACKAGE/$GOFILE . CognitoClient
package aws

import (
	"context"
	"crypto/hmac"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	appConfig "share-basket-server/core/config"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/golang-jwt/jwt/v5"
)

type (
	CognitoClient interface {
		InitiateAuth(ctx context.Context, email, password string) (*cognitoidentityprovider.InitiateAuthOutput, error)
		SignUp(ctx context.Context, email, password string) (*cognitoidentityprovider.SignUpOutput, error)
		ConfirmSignUp(ctx context.Context, email, confirmationCode string) (*cognitoidentityprovider.ConfirmSignUpOutput, error)
		AdminDeleteUser(ctx context.Context, email string) error
		ParseToken(ctx context.Context, token string) (*jwt.Token, error)
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

func (client *cognitoClient) ParseToken(ctx context.Context, token string) (*jwt.Token, error) {
	return jwt.Parse(token, client.keyFunc)
}

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

// JWTの`kid`に対応する公開鍵を取得する
func (client *cognitoClient) keyFunc(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}

	kid, ok := token.Header["kid"].(string)
	if !ok {
		return nil, errors.New("missing kid in token header")
	}

	client.mu.Lock()
	defer client.mu.Unlock()

	if key, exists := client.keys[kid]; exists {
		return key, nil
	}

	key, err := client.fetchJWKS(kid)
	if err != nil {
		return nil, err
	}

	client.keys[kid] = key
	return key, nil
}

// CognitoのJWKSエンドポイントから公開鍵を取得する
func (client *cognitoClient) fetchJWKS(kid string) (*rsa.PublicKey, error) {
	resp, err := http.Get(client.jwksURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch JWKS: %w", err)
	}
	defer resp.Body.Close()

	var jwks struct {
		Keys []struct {
			Kty string `json:"kty"`
			Kid string `json:"kid"`
			N   string `json:"n"`
			E   string `json:"e"`
		} `json:"keys"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return nil, fmt.Errorf("failed to decode JWKS: %w", err)
	}

	for _, key := range jwks.Keys {
		if key.Kid == kid {
			nBytes, err := base64.RawURLEncoding.DecodeString(key.N)
			if err != nil {
				return nil, fmt.Errorf("failed to decode N: %w", err)
			}

			eBytes, err := base64.RawURLEncoding.DecodeString(key.E)
			if err != nil {
				return nil, fmt.Errorf("failed to decode E: %w", err)
			}

			e := 0
			for _, b := range eBytes {
				e = e<<8 + int(b)
			}

			rsaKey := &rsa.PublicKey{
				N: new(big.Int).SetBytes(nBytes),
				E: e,
			}

			return rsaKey, nil
		}
	}

	return nil, errors.New("public key not found in JWKS")
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
