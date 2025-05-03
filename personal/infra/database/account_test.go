package database_test

import (
	"share-basket-server/personal/domain"
	"share-basket-server/personal/infra/database"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccountPersistence(t *testing.T) {
	repo := database.NewAccountPersistence(testDB)

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
		userDto := database.UserDto{
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
