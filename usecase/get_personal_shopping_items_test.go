package usecase_test

import (
	"context"
	"errors"
	"testing"

	"share-basket-server/core/apperr"
	contextKey "share-basket-server/core/context"
	"share-basket-server/core/util"
	"share-basket-server/domain"
	"share-basket-server/usecase"

	. "share-basket-server/mock/domain"
	. "share-basket-server/mock/usecase"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestGetPersonalShoppingItemsInteractor(t *testing.T) {
	var (
		userID       = domain.NewUserID()
		ctx          = context.WithValue(context.Background(), contextKey.UserID, domain.UserID("dummy-user-id"))
		ctrl         = gomock.NewController(t)
		accountRepo  = NewMockAccountRepository(ctrl)
		personalRepo = NewMockPersonalShoppingItemRepository(ctrl)
		output       = NewMockGetPersonalShoppingItemsOutputPort(ctrl)
	)

	defer ctrl.Finish()

	tests := map[string]struct {
		setupMock func(
			accRepo *MockAccountRepository,
			personalRepo *MockPersonalShoppingItemRepository,
			output *MockGetPersonalShoppingItemsOutputPort,
		)
		input usecase.GetPersonalShoppingItemsInput
		err   error
	}{
		"正常系: 買い物リストの取得に成功": {
			setupMock: func(
				accRepo *MockAccountRepository,
				personalRepo *MockPersonalShoppingItemRepository,
				output *MockGetPersonalShoppingItemsOutputPort,
			) {
				account := domain.Account{
					ID:     domain.AccountID("dummy-account-id"),
					UserID: userID,
					Name:   "dummy user",
				}
				accRepo.EXPECT().
					FindByUserID(userID).
					Return(account, nil)

				items := []domain.PersonalShoppingItem{
					{
						ID:         util.Ptr[uint](1),
						Name:       "牛乳",
						Status:     domain.UnPurchased,
						CategoryID: 1,
						AccountID:  domain.AccountID("dummy-account-id"),
					},
				}
				personalRepo.EXPECT().
					GetAll(account.ID, nil).
					Return(items, nil)

				output.EXPECT().
					Render(ctx, []usecase.GetPersonalShoppingItemOutput{
						{
							ID:         1,
							Name:       "牛乳",
							CategoryID: 1,
							Status:     "UnPurchased",
						},
					}).
					Return(nil)
			},
			input: usecase.GetPersonalShoppingItemsInput{
				UserID: userID.String(),
				Status: "",
			},
			err: nil,
		},
		"正常系: ステータスを指定して買い物リストの取得に成功": {
			setupMock: func(
				accRepo *MockAccountRepository,
				personalRepo *MockPersonalShoppingItemRepository,
				output *MockGetPersonalShoppingItemsOutputPort,
			) {
				account := domain.Account{
					ID:     domain.AccountID("dummy-account-id"),
					UserID: userID,
					Name:   "dummy user",
				}
				accRepo.EXPECT().
					FindByUserID(userID).
					Return(account, nil)

				items := []domain.PersonalShoppingItem{
					{
						ID:         util.Ptr[uint](1),
						Name:       "牛乳",
						Status:     domain.UnPurchased,
						CategoryID: 1,
						AccountID:  domain.AccountID("dummy-account-id"),
					},
				}
				personalRepo.EXPECT().
					GetAll(account.ID, util.Ptr(domain.UnPurchased)).
					Return(items, nil)

				output.EXPECT().
					Render(ctx, []usecase.GetPersonalShoppingItemOutput{
						{
							ID:         1,
							Name:       "牛乳",
							CategoryID: 1,
							Status:     "UnPurchased",
						},
					}).
					Return(nil)
			},
			input: usecase.GetPersonalShoppingItemsInput{
				UserID: userID.String(),
				Status: domain.UnPurchased.String(),
			},
			err: nil,
		},
		"異常系: UserIDが正しくない": {
			setupMock: func(accRepo *MockAccountRepository, personalRepo *MockPersonalShoppingItemRepository, output *MockGetPersonalShoppingItemsOutputPort) {
			},
			input: usecase.GetPersonalShoppingItemsInput{
				UserID: "invalid-id",
				Status: "",
			},
			err: apperr.NewInvalidError(errors.New("invalid user id: invalid UUID length: 10")),
		},
		"異常系: アカウントの取得に失敗": {
			setupMock: func(
				accRepo *MockAccountRepository,
				personalRepo *MockPersonalShoppingItemRepository,
				output *MockGetPersonalShoppingItemsOutputPort,
			) {
				accRepo.EXPECT().
					FindByUserID(userID).
					Return(domain.Account{}, apperr.ErrDataNotFound)
			},
			input: usecase.GetPersonalShoppingItemsInput{
				UserID: userID.String(),
				Status: "",
			},
			err: apperr.NewInvalidError(apperr.ErrDataNotFound),
		},
		"異常系: 不正なステータスを指定": {
			setupMock: func(
				accRepo *MockAccountRepository,
				personalRepo *MockPersonalShoppingItemRepository,
				output *MockGetPersonalShoppingItemsOutputPort,
			) {
				account := domain.Account{
					ID:     domain.AccountID("dummy-account-id"),
					UserID: userID,
					Name:   "dummy user",
				}
				accRepo.EXPECT().
					FindByUserID(userID).
					Return(account, nil)
			},
			input: usecase.GetPersonalShoppingItemsInput{
				UserID: userID.String(),
				Status: "InvalidStatus",
			},
			err: apperr.NewInvalidError(errors.New("invalid shopping status")),
		},
		"異常系: 買い物リストの取得に失敗": {
			setupMock: func(
				accRepo *MockAccountRepository,
				personalRepo *MockPersonalShoppingItemRepository,
				output *MockGetPersonalShoppingItemsOutputPort,
			) {
				account := domain.Account{
					ID:     domain.AccountID("dummy-account-id"),
					UserID: userID,
					Name:   "dummy user",
				}
				accRepo.EXPECT().
					FindByUserID(userID).
					Return(account, nil)

				personalRepo.EXPECT().
					GetAll(account.ID, nil).
					Return(nil, errors.New("database error"))
			},
			input: usecase.GetPersonalShoppingItemsInput{
				UserID: userID.String(),
				Status: "",
			},
			err: errors.New("database error"),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			tt.setupMock(accountRepo, personalRepo, output)

			usecase := usecase.NewGetPersonalShoppingItemsInteractor(accountRepo, personalRepo)
			err := usecase.Execute(ctx, tt.input, output)

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
