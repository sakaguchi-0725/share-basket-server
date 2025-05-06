package handler_test

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"testing"

	"share-basket-server/core/apperr"
	. "share-basket-server/mock/usecase"
	. "share-basket-server/mock/validator"
	"share-basket-server/presentation/handler"
	"share-basket-server/usecase"

	"go.uber.org/mock/gomock"
	"gopkg.in/go-playground/assert.v1"
)

func TestLoginHandler(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()

	t.Run("MakeLoginHandler", func(t *testing.T) {
		var (
			server     = newTestServer(http.MethodPost, "/login")
			interactor = NewMockLoginInputPort(ctrl)
			validator  = NewMockRequestValidator(ctrl)
		)

		tests := map[string]struct {
			setupMock  func(validator *MockRequestValidator, interactor *MockLoginInputPort)
			reqBody    string
			wantStatus int
		}{
			"正常系: ログインに成功した場合": {
				setupMock: func(validator *MockRequestValidator, interactor *MockLoginInputPort) {
					validator.EXPECT().Validate(&handler.LoginRequest{
						Email:    "test@example.com",
						Password: "password",
					}).Return(nil)
					interactor.EXPECT().Execute(
						gomock.Any(),
						usecase.LoginInput{
							Email:    "test@example.com",
							Password: "password",
						},
						gomock.Any(),
					).DoAndReturn(func(_ context.Context, _ usecase.LoginInput, out usecase.LoginOutputPort) error {
						return out.Render(context.Background(), "dummy-token")
					})
				},
				reqBody:    `{"email":"test@example.com","password":"password"}`,
				wantStatus: http.StatusNoContent,
			},
			"異常系: メールアドレス形式が不正": {
				setupMock: func(validator *MockRequestValidator, interactor *MockLoginInputPort) {
					validator.EXPECT().Validate(&handler.LoginRequest{
						Email:    "invalid-email",
						Password: "password",
					}).Return(errors.New("invalid email format"))
				},
				reqBody:    `{"email":"invalid-email","password":"password"}`,
				wantStatus: http.StatusBadRequest,
			},
			"異常系: 必須項目が未入力": {
				setupMock: func(validator *MockRequestValidator, interactor *MockLoginInputPort) {
					validator.EXPECT().Validate(&handler.LoginRequest{
						Email:    "",
						Password: "",
					}).Return(errors.New("email and password are required"))
				},
				reqBody:    `{"email":"","password":""}`,
				wantStatus: http.StatusBadRequest,
			},
			"異常系: JSONデコードエラー": {
				setupMock: func(validator *MockRequestValidator, interactor *MockLoginInputPort) {
					// バリデーションは呼ばれない
				},
				reqBody:    `{"email":"test@example.com","password":"password"`, // 閉じ括弧がない
				wantStatus: http.StatusBadRequest,
			},
			"異常系: ログインに失敗した場合": {
				setupMock: func(validator *MockRequestValidator, interactor *MockLoginInputPort) {
					validator.EXPECT().Validate(&handler.LoginRequest{
						Email:    "notfound@example.com",
						Password: "password",
					}).Return(nil)
					interactor.EXPECT().Execute(
						gomock.Any(),
						usecase.LoginInput{
							Email:    "notfound@example.com",
							Password: "password",
						},
						gomock.Any(),
					).Return(apperr.NewInvalidError(errors.New("user not found")))
				},
				reqBody:    `{"email":"notfound@example.com","password":"password"}`,
				wantStatus: http.StatusBadRequest,
			},
		}

		for name, tt := range tests {
			t.Run(name, func(t *testing.T) {
				tt.setupMock(validator, interactor)

				handler := handler.MakeLoginHandler(interactor, validator)
				rec := server.Serve(bytes.NewReader([]byte(tt.reqBody)), handler)
				assert.Equal(t, tt.wantStatus, rec.Code)
			})
		}
	})
}
