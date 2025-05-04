package usecase_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"share-basket-server/core/apperr"
	. "share-basket-server/mock/domain"
	. "share-basket-server/mock/usecase"
	"share-basket-server/usecase"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestSignUpInteractor(t *testing.T) {
	var (
		ctrl          = gomock.NewController(t)
		authenticator = NewMockAuthenticator(ctrl)
		userRepo      = NewMockUserRepository(ctrl)
		accountRepo   = NewMockAccountRepository(ctrl)
		userService   = NewMockUserService(ctrl)
		transaction   = NewMockTransaction(ctrl)
		output        = NewMockSignUpOutputPort(ctrl)
	)

	defer ctrl.Finish()

	tests := map[string]struct {
		setupMock func(
			authenticator *MockAuthenticator,
			userRepo *MockUserRepository,
			accountRepo *MockAccountRepository,
			userService *MockUserService,
			transaction *MockTransaction,
			output *MockSignUpOutputPort,
		)
		input usecase.SignUpInput
		err   error
	}{
		"正常系: サインアップ成功": {
			setupMock: func(
				authenticator *MockAuthenticator,
				userRepo *MockUserRepository,
				accountRepo *MockAccountRepository,
				userService *MockUserService,
				transaction *MockTransaction,
				output *MockSignUpOutputPort,
			) {
				userService.EXPECT().IsEmailAvailable("test@example.com").Return(true, nil)
				authenticator.EXPECT().SignUp(context.Background(), "test@example.com", "password").Return("cognito-uid", nil)
				transaction.EXPECT().Run(context.Background(), gomock.Any()).DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
					return fn(ctx)
				})
				userRepo.EXPECT().Store(gomock.Any()).Return(nil)
				accountRepo.EXPECT().Store(gomock.Any()).Return(nil)
				output.EXPECT().Render(context.Background()).Return(nil)
			},
			input: usecase.SignUpInput{
				Name:     "Test User",
				Email:    "test@example.com",
				Password: "password",
			},
			err: nil,
		},
		"異常系: メールアドレスが利用不可": {
			setupMock: func(
				authenticator *MockAuthenticator,
				userRepo *MockUserRepository,
				accountRepo *MockAccountRepository,
				userService *MockUserService,
				transaction *MockTransaction,
				output *MockSignUpOutputPort,
			) {
				userService.EXPECT().IsEmailAvailable("test@example.com").Return(false, nil)
			},
			input: usecase.SignUpInput{
				Name:     "Test User",
				Email:    "test@example.com",
				Password: "password",
			},
			err: apperr.NewInvalidError(fmt.Errorf("email is not available: %s", "test@example.com")),
		},
		"異常系: メールアドレスチェックでエラー": {
			setupMock: func(
				authenticator *MockAuthenticator,
				userRepo *MockUserRepository,
				accountRepo *MockAccountRepository,
				userService *MockUserService,
				transaction *MockTransaction,
				output *MockSignUpOutputPort,
			) {
				userService.EXPECT().IsEmailAvailable("test@example.com").Return(false, errors.New("check error"))
			},
			input: usecase.SignUpInput{
				Name:     "Test User",
				Email:    "test@example.com",
				Password: "password",
			},
			err: errors.New("check error"),
		},
		"異常系: サインアップで重複エラー": {
			setupMock: func(
				authenticator *MockAuthenticator,
				userRepo *MockUserRepository,
				accountRepo *MockAccountRepository,
				userService *MockUserService,
				transaction *MockTransaction,
				output *MockSignUpOutputPort,
			) {
				userService.EXPECT().IsEmailAvailable("test@example.com").Return(true, nil)
				authenticator.EXPECT().SignUp(context.Background(), "test@example.com", "password").Return("", apperr.ErrDuplicatedKey)
			},
			input: usecase.SignUpInput{
				Name:     "Test User",
				Email:    "test@example.com",
				Password: "password",
			},
			err: apperr.NewInvalidError(apperr.ErrDuplicatedKey),
		},
		"異常系: トランザクションでエラー": {
			setupMock: func(
				authenticator *MockAuthenticator,
				userRepo *MockUserRepository,
				accountRepo *MockAccountRepository,
				userService *MockUserService,
				transaction *MockTransaction,
				output *MockSignUpOutputPort,
			) {
				userService.EXPECT().IsEmailAvailable("test@example.com").Return(true, nil)
				authenticator.EXPECT().SignUp(context.Background(), "test@example.com", "password").Return("cognito-uid", nil)
				transaction.EXPECT().Run(context.Background(), gomock.Any()).Return(errors.New("transaction error"))
				authenticator.EXPECT().DeleteUser(context.Background(), "test@example.com").Return(nil)
			},
			input: usecase.SignUpInput{
				Name:     "Test User",
				Email:    "test@example.com",
				Password: "password",
			},
			err: errors.New("transaction error"),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			tt.setupMock(authenticator, userRepo, accountRepo, userService, transaction, output)

			usecase := usecase.NewSignUpInteractor(authenticator, userRepo, accountRepo, userService, transaction)
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
