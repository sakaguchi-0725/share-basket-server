package usecase_test

import (
	"context"
	"errors"
	"testing"

	"share-basket-server/core/apperr"
	. "share-basket-server/test/mock/domain"
	. "share-basket-server/test/mock/usecase"
	"share-basket-server/test/testutil"
	"share-basket-server/usecase"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestLoginInteractor(t *testing.T) {
	var (
		ctrl          = gomock.NewController(t)
		authenticator = NewMockAuthenticator(ctrl)
		output        = NewMockLoginOutputPort(ctrl)
	)

	defer ctrl.Finish()

	tests := map[string]struct {
		setupMock func(authenticator *MockAuthenticator, output *MockLoginOutputPort)
		input     usecase.LoginInput
		err       error
		want      string
	}{
		"正常系: ログイン成功": {
			setupMock: func(authenticator *MockAuthenticator, output *MockLoginOutputPort) {
				authenticator.EXPECT().Login(context.Background(), "test@example.com", "password").Return("access_token", nil)
				output.EXPECT().Render(context.Background(), "access_token").Return(nil)
			},
			input: usecase.LoginInput{
				Email:    "test@example.com",
				Password: "password",
			},
			err:  nil,
			want: "access_token",
		},
		"異常系: 認証エラー": {
			setupMock: func(authenticator *MockAuthenticator, output *MockLoginOutputPort) {
				authenticator.EXPECT().Login(context.Background(), "test@example.com", "wrong_password").Return("", apperr.ErrUnauthenticated)
			},
			input: usecase.LoginInput{
				Email:    "test@example.com",
				Password: "wrong_password",
			},
			err:  apperr.New(apperr.ErrUnauthorized, apperr.ErrUnauthenticated),
			want: "",
		},
		"異常系: ユーザーが見つからない": {
			setupMock: func(authenticator *MockAuthenticator, output *MockLoginOutputPort) {
				authenticator.EXPECT().Login(context.Background(), "notfound@example.com", "password").Return("", apperr.ErrDataNotFound)
			},
			input: usecase.LoginInput{
				Email:    "notfound@example.com",
				Password: "password",
			},
			err:  apperr.NewInvalidError(apperr.ErrDataNotFound),
			want: "",
		},
		"異常系: その他のエラー": {
			setupMock: func(authenticator *MockAuthenticator, output *MockLoginOutputPort) {
				authenticator.EXPECT().Login(context.Background(), "test@example.com", "password").Return("", errors.New("internal error"))
			},
			input: usecase.LoginInput{
				Email:    "test@example.com",
				Password: "password",
			},
			err:  errors.New("internal error"),
			want: "",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			tt.setupMock(authenticator, output)

			usecase := usecase.NewLoginInteractor(authenticator, testutil.NewDummyLogger())
			err := usecase.Execute(context.Background(), tt.input, output)

			if tt.err == nil {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Equal(t, tt.err.Error(), err.Error())

				if appErr, ok := tt.err.(*apperr.AppError); ok {
					if gotAppErr, ok := err.(*apperr.AppError); ok {
						assert.Equal(t, appErr.Code(), gotAppErr.Code())
					}
				}
			}
		})
	}
}
