package auth

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
	"sharebasket/core"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
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
		client     *cognitoidentityprovider.Client
		id         string
		secret     string
		userPoolID string
		jwksURL    string
		keys       map[string]*rsa.PublicKey
		mu         sync.Mutex
	}
)

func (c *cognitoClient) ParseToken(ctx context.Context, token string) (*jwt.Token, error) {
	return jwt.Parse(token, c.keyFunc)
}

func (c *cognitoClient) AdminDeleteUser(ctx context.Context, email string) error {
	input := &cognitoidentityprovider.AdminDeleteUserInput{
		UserPoolId: aws.String(c.userPoolID),
		Username:   aws.String(email),
	}

	_, err := c.client.AdminDeleteUser(ctx, input)
	if err != nil {
		var notFound *types.UserNotFoundException
		if errors.As(err, &notFound) {
			return nil // ユーザーが存在しない場合は成功として扱う
		}
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

func (c *cognitoClient) ConfirmSignUp(ctx context.Context, email, confirmationCode string) (*cognitoidentityprovider.ConfirmSignUpOutput, error) {
	return c.client.ConfirmSignUp(ctx, &cognitoidentityprovider.ConfirmSignUpInput{
		ClientId:         aws.String(c.id),
		Username:         aws.String(email),
		ConfirmationCode: aws.String(confirmationCode),
		SecretHash:       aws.String(c.genSecretHash(email)),
	})
}

func (c *cognitoClient) InitiateAuth(ctx context.Context, email, password string) (*cognitoidentityprovider.InitiateAuthOutput, error) {
	return c.client.InitiateAuth(ctx, &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: types.AuthFlowTypeUserPasswordAuth,
		ClientId: aws.String(c.id),
		AuthParameters: map[string]string{
			"USERNAME":    email,
			"PASSWORD":    password,
			"SECRET_HASH": c.genSecretHash(email),
		},
	})
}

func (c *cognitoClient) SignUp(ctx context.Context, email, password string) (*cognitoidentityprovider.SignUpOutput, error) {
	return c.client.SignUp(ctx, &cognitoidentityprovider.SignUpInput{
		ClientId:   aws.String(c.id),
		Username:   aws.String(email),
		Password:   aws.String(password),
		SecretHash: aws.String(c.genSecretHash(email)),
	})
}

// Cognito認証に必要なSecretHashを生成する
func (c *cognitoClient) genSecretHash(username string) string {
	mac := hmac.New(sha256.New, []byte(c.secret))
	mac.Write([]byte(username + c.id))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

// JWTの`kid`に対応する公開鍵を取得する
func (c *cognitoClient) keyFunc(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}

	kid, ok := token.Header["kid"].(string)
	if !ok {
		return nil, errors.New("missing kid in token header")
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if key, exists := c.keys[kid]; exists {
		return key, nil
	}

	key, err := c.fetchJWKS(kid)
	if err != nil {
		return nil, err
	}

	c.keys[kid] = key
	return key, nil
}

// CognitoのJWKSエンドポイントから公開鍵を取得する
func (c *cognitoClient) fetchJWKS(kid string) (*rsa.PublicKey, error) {
	resp, err := http.Get(c.jwksURL)
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

func NewCognitoClient(ctx context.Context) (CognitoClient, error) {
	region := core.GetEnv("AWS_REGION", "ap-northeast-1")
	userPoolID := core.GetEnv("COGNITO_USER_POOL_ID", "")

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				core.GetEnv("AWS_ACCESS_KEY_ID", ""),
				core.GetEnv("AWS_SECRET_ACCESS_KEY", ""),
				"",
			),
		),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	return &cognitoClient{
		client:     cognitoidentityprovider.NewFromConfig(cfg),
		id:         core.GetEnv("COGNITO_CLIENT_ID", ""),
		secret:     core.GetEnv("COGNITO_CLIENT_SECRET", ""),
		userPoolID: userPoolID,
		jwksURL: fmt.Sprintf(
			"https://cognito-idp.%s.amazonaws.com/%s/.well-known/jwks.json",
			region, userPoolID,
		),
		keys: make(map[string]*rsa.PublicKey),
	}, nil
}
