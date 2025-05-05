package repository_test

import (
	"testing"

	"share-basket-server/core/apperr"
	"share-basket-server/domain"
	"share-basket-server/infra/rdb/repository"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccountPersistence(t *testing.T) {
	repo := repository.NewAccountPersistence(testDB)

	t.Run("FindByUserID", func(t *testing.T) {
		defer clearTestData()

		err := createDummyAccount()
		require.NoError(t, err)

		tests := map[string]struct {
			userID domain.UserID
			want   domain.Account
			err    error
		}{
			"正常系: アカウントが取得できる": {
				userID: domain.UserID(dummyUser.ID),
				want: domain.Account{
					ID:     domain.AccountID(dummyAccount.ID),
					UserID: domain.UserID(dummyAccount.UserID),
					Name:   dummyAccount.Name,
				},
				err: nil,
			},
			"異常系: UserIDに紐づくアカウントが存在しない": {
				userID: domain.NewUserID(),
				want:   domain.Account{},
				err:    apperr.ErrDataNotFound,
			},
		}

		for name, tt := range tests {
			t.Run(name, func(t *testing.T) {
				got, err := repo.FindByUserID(tt.userID)

				assert.Equal(t, tt.want, got)
				if tt.err != nil {
					assert.Error(t, err)
					assert.EqualError(t, tt.err, err.Error())
				} else {
					assert.NoError(t, err)
				}
			})
		}
	})

	t.Run("Store", func(t *testing.T) {
		defer clearTestData()

		var (
			accID  = domain.NewAccountID()
			userID = domain.NewUserID()
			user   = domain.User{
				ID:         userID,
				CognitoUID: "cognito-uid",
				Email:      "test@example.com",
			}
		)

		// 事前にユーザーを作成
		userDto := repository.UserDto{
			ID:         user.ID.String(),
			CognitoUID: user.CognitoUID,
			Email:      user.Email,
		}
		err := testDB.Create(&userDto).Error
		require.NoError(t, err, "ユーザーの作成に失敗しました")

		tests := map[string]struct {
			input *domain.Account
			err   error
		}{
			"正常系: アカウントが作成される": {
				input: &domain.Account{
					ID:     accID,
					UserID: userID,
					Name:   "test",
				},
				err: nil,
			},
		}

		for name, tt := range tests {
			t.Run(name, func(t *testing.T) {
				err := repo.Store(tt.input)
				if tt.err != nil {
					assert.Error(t, err)
					assert.Equal(t, tt.err, err)
					return
				}

				assert.NoError(t, err)
			})
		}
	})
}
