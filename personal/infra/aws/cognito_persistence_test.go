package aws_test

import (
	"context"
	"errors"
	"share-basket-server/core/apperr"
	"share-basket-server/personal/infra/aws"
	mock_aws "share-basket-server/personal/mock/aws"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCognitoPersistence(t *testing.T) {
	var (
		ctrl          = gomock.NewController(t)
		ctx           = context.Background()
		mockClient    = mock_aws.NewMockCognitoClient(ctrl)
		authenticator = aws.NewCognitoPersistence(mockClient)
	)

	defer ctrl.Finish()

	t.Run("Login", func(t *testing.T) {
		cases := map[string]struct {
			email    string
			password string
			setup    func()
			want     string
			wantErr  error
		}{
			"正常系: ログイン成功": {
				email:    "test@example.com",
				password: "password123",
				setup: func() {
					accessToken := "test-access-token"
					mockClient.EXPECT().InitiateAuth(ctx, "test@example.com", "password123").Return(&cognitoidentityprovider.InitiateAuthOutput{
						AuthenticationResult: &types.AuthenticationResultType{
							AccessToken: &accessToken,
						},
					}, nil)
				},
				want: "test-access-token",
			},
			"異常系: 認証エラー": {
				email:    "test@example.com",
				password: "wrong-password",
				setup: func() {
					mockClient.EXPECT().InitiateAuth(ctx, "test@example.com", "wrong-password").Return(nil, &types.NotAuthorizedException{})
				},
				wantErr: apperr.ErrUnauthenticated,
			},
			"異常系: 無効なパラメータ": {
				email:    "test@example.com",
				password: "",
				setup: func() {
					mockClient.EXPECT().InitiateAuth(ctx, "test@example.com", "").Return(nil, &types.InvalidParameterException{})
				},
				wantErr: apperr.ErrInvalidData,
			},
		}

		for name, tc := range cases {
			t.Run(name, func(t *testing.T) {
				tc.setup()
				got, err := authenticator.Login(ctx, tc.email, tc.password)
				if tc.wantErr != nil {
					assert.ErrorIs(t, err, tc.wantErr)
					return
				}
				assert.NoError(t, err)
				assert.Equal(t, tc.want, got)
			})
		}
	})

	t.Run("SignUp", func(t *testing.T) {
		cases := map[string]struct {
			email    string
			password string
			setup    func()
			want     string
			wantErr  error
		}{
			"正常系: サインアップ成功": {
				email:    "test@example.com",
				password: "password123",
				setup: func() {
					userSub := "test-user-sub"
					mockClient.EXPECT().SignUp(ctx, "test@example.com", "password123").Return(&cognitoidentityprovider.SignUpOutput{
						UserSub: &userSub,
					}, nil)
				},
				want: "test-user-sub",
			},
			"異常系: ユーザーが既に存在": {
				email:    "test@example.com",
				password: "password123",
				setup: func() {
					mockClient.EXPECT().SignUp(ctx, "test@example.com", "password123").Return(nil, &types.UsernameExistsException{})
				},
				wantErr: apperr.ErrDuplicatedKey,
			},
			"異常系: 無効なパラメータ": {
				email:    "test@example.com",
				password: "",
				setup: func() {
					mockClient.EXPECT().SignUp(ctx, "test@example.com", "").Return(nil, &types.InvalidParameterException{})
				},
				wantErr: apperr.ErrInvalidData,
			},
		}

		for name, tc := range cases {
			t.Run(name, func(t *testing.T) {
				tc.setup()
				got, err := authenticator.SignUp(ctx, tc.email, tc.password)
				if tc.wantErr != nil {
					assert.ErrorIs(t, err, tc.wantErr)
					return
				}
				assert.NoError(t, err)
				assert.Equal(t, tc.want, got)
			})
		}
	})

	t.Run("SignUpConfirm", func(t *testing.T) {
		cases := map[string]struct {
			email string
			code  string
			setup func()
			want  error
		}{
			"正常系: サインアップ確認成功": {
				email: "test@example.com",
				code:  "123456",
				setup: func() {
					mockClient.EXPECT().ConfirmSignUp(ctx, "test@example.com", "123456").Return(&cognitoidentityprovider.ConfirmSignUpOutput{}, nil)
				},
			},
			"異常系: 無効な確認コード": {
				email: "test@example.com",
				code:  "wrong-code",
				setup: func() {
					mockClient.EXPECT().ConfirmSignUp(ctx, "test@example.com", "wrong-code").Return(nil, &types.CodeMismatchException{})
				},
				want: apperr.ErrInvalidData,
			},
			"異常系: 期限切れの確認コード": {
				email: "test@example.com",
				code:  "expired-code",
				setup: func() {
					mockClient.EXPECT().ConfirmSignUp(ctx, "test@example.com", "expired-code").Return(nil, &types.ExpiredCodeException{})
				},
				want: apperr.ErrExpiredCodeException,
			},
		}

		for name, tc := range cases {
			t.Run(name, func(t *testing.T) {
				tc.setup()
				err := authenticator.SignUpConfirm(ctx, tc.email, tc.code)
				if tc.want != nil {
					assert.ErrorIs(t, err, tc.want)
					return
				}
				assert.NoError(t, err)
			})
		}
	})

	t.Run("DeleteUser", func(t *testing.T) {
		cases := map[string]struct {
			email string
			setup func()
			want  error
		}{
			"正常系: ユーザー削除成功": {
				email: "test@example.com",
				setup: func() {
					mockClient.EXPECT().AdminDeleteUser(ctx, "test@example.com").Return(nil)
				},
			},
			"異常系: ユーザー削除失敗": {
				email: "test@example.com",
				setup: func() {
					mockClient.EXPECT().AdminDeleteUser(ctx, "test@example.com").Return(errors.New("delete failed"))
				},
				want: errors.New("delete failed"),
			},
		}

		for name, tc := range cases {
			t.Run(name, func(t *testing.T) {
				tc.setup()
				err := authenticator.DeleteUser(ctx, tc.email)
				if tc.want != nil {
					assert.Error(t, err)
					return
				}
				assert.NoError(t, err)
			})
		}
	})
}
