package usecase_test

import (
	"context"
	"errors"
	"testing"

	"share-basket-server/core/apperr"
	"share-basket-server/domain"
	. "share-basket-server/test/mock/domain"
	"share-basket-server/test/testutil"
	"share-basket-server/usecase"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestVerifyTokenInteractor(t *testing.T) {
	var (
		ctrl          = gomock.NewController(t)
		authenticator = NewMockAuthenticator(ctrl)
		userRepo      = NewMockUserRepository(ctrl)
	)
	defer ctrl.Finish()

	userID := domain.NewUserID()

	tests := map[string]struct {
		setupMock func(authenticator *MockAuthenticator, userRepo *MockUserRepository)
		token     string
		want      string
		err       error
	}{
		"正常系: トークン検証成功": {
			setupMock: func(authenticator *MockAuthenticator, userRepo *MockUserRepository) {
				authenticator.EXPECT().VerifyToken(context.Background(), "valid-token").Return("test@example.com", nil)
				userRepo.EXPECT().GetByEmail("test@example.com").Return(domain.User{ID: userID}, nil)
			},
			token: "valid-token",
			want:  userID.String(),
			err:   nil,
		},
		"異常系: 無効なトークン": {
			setupMock: func(authenticator *MockAuthenticator, userRepo *MockUserRepository) {
				authenticator.EXPECT().VerifyToken(context.Background(), "invalid-token").Return("", apperr.ErrInvalidToken)
			},
			token: "invalid-token",
			want:  "",
			err:   apperr.New(apperr.ErrUnauthorized, apperr.ErrInvalidToken),
		},
		"異常系: トークンの有効期限切れ": {
			setupMock: func(authenticator *MockAuthenticator, userRepo *MockUserRepository) {
				authenticator.EXPECT().VerifyToken(context.Background(), "expired-token").Return("", apperr.ErrTokenExpired)
			},
			token: "expired-token",
			want:  "",
			err:   apperr.New(apperr.ErrUnauthorized, apperr.ErrTokenExpired),
		},
		"異常系: ユーザーが見つからない": {
			setupMock: func(authenticator *MockAuthenticator, userRepo *MockUserRepository) {
				authenticator.EXPECT().VerifyToken(context.Background(), "valid-token").Return("test@example.com", nil)
				userRepo.EXPECT().GetByEmail("test@example.com").Return(domain.User{}, apperr.ErrDataNotFound)
			},
			token: "valid-token",
			want:  "",
			err:   apperr.New(apperr.ErrUnauthorized, apperr.ErrDataNotFound),
		},
		"異常系: その他のエラー": {
			setupMock: func(authenticator *MockAuthenticator, userRepo *MockUserRepository) {
				authenticator.EXPECT().VerifyToken(context.Background(), "error-token").Return("", errors.New("internal error"))
			},
			token: "error-token",
			want:  "",
			err:   errors.New("internal error"),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			tt.setupMock(authenticator, userRepo)

			usecase := usecase.NewVerifyTokenInteractor(authenticator, userRepo, testutil.NewDummyLogger())
			got, err := usecase.Execute(context.Background(), tt.token)

			if tt.err == nil {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			} else {
				assert.Error(t, err)
				assert.Equal(t, tt.err.Error(), err.Error())
				assert.Equal(t, tt.want, got)

				if appErr, ok := tt.err.(*apperr.AppError); ok {
					if gotAppErr, ok := err.(*apperr.AppError); ok {
						assert.Equal(t, appErr.Code(), gotAppErr.Code())
					}
				}
			}
		})
	}
}
