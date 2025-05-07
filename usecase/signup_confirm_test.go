package usecase_test

import (
	"context"
	"testing"

	"share-basket-server/core/apperr"
	. "share-basket-server/test/mock/domain"
	. "share-basket-server/test/mock/usecase"
	"share-basket-server/usecase"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestSignUpConfirmInteractor(t *testing.T) {
	var (
		ctrl          = gomock.NewController(t)
		authenticator = NewMockAuthenticator(ctrl)
		output        = NewMockSignUpConfirmOutputPort(ctrl)
	)

	defer ctrl.Finish()

	tests := map[string]struct {
		setupMock func(authenticator *MockAuthenticator, output *MockSignUpConfirmOutputPort)
		input     usecase.SignUpConfirmInput
		err       error
	}{
		"正常系: サインアップ確認成功": {
			setupMock: func(authenticator *MockAuthenticator, output *MockSignUpConfirmOutputPort) {
				authenticator.EXPECT().SignUpConfirm(context.Background(), "test@example.com", "123456").Return(nil)
				output.EXPECT().Render(context.Background()).Return(nil)
			},
			input: usecase.SignUpConfirmInput{
				Email:            "test@example.com",
				ConfirmationCode: "123456",
			},
			err: nil,
		},
		"異常系: 無効なデータ": {
			setupMock: func(authenticator *MockAuthenticator, output *MockSignUpConfirmOutputPort) {
				authenticator.EXPECT().SignUpConfirm(context.Background(), "test@example.com", "invalid").Return(apperr.ErrInvalidData)
			},
			input: usecase.SignUpConfirmInput{
				Email:            "test@example.com",
				ConfirmationCode: "invalid",
			},
			err: apperr.NewInvalidError(apperr.ErrInvalidData),
		},
		"異常系: コードの有効期限切れ": {
			setupMock: func(authenticator *MockAuthenticator, output *MockSignUpConfirmOutputPort) {
				authenticator.EXPECT().SignUpConfirm(context.Background(), "test@example.com", "expired").Return(apperr.ErrExpiredCodeException)
			},
			input: usecase.SignUpConfirmInput{
				Email:            "test@example.com",
				ConfirmationCode: "expired",
			},
			err: apperr.New(apperr.ErrExpiredCode, apperr.ErrExpiredCodeException),
		},
		"異常系: その他のエラー": {
			setupMock: func(authenticator *MockAuthenticator, output *MockSignUpConfirmOutputPort) {
				authenticator.EXPECT().SignUpConfirm(context.Background(), "test@example.com", "error").Return(apperr.ErrDataNotFound)
			},
			input: usecase.SignUpConfirmInput{
				Email:            "test@example.com",
				ConfirmationCode: "error",
			},
			err: apperr.ErrDataNotFound,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			tt.setupMock(authenticator, output)

			usecase := usecase.NewSignUpConfirmInteractor(authenticator)
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
