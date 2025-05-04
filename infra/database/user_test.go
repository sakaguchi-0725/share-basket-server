package database_test

import (
	"share-basket-server/core/apperr"
	"share-basket-server/domain"
	"share-basket-server/infra/database"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserPersistence(t *testing.T) {
	repo := database.NewUserPersistence(testDB)

	t.Run("GetByEmail", func(t *testing.T) {
		defer clearTestData()

		var (
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
			email string
			want  domain.User
			err   error
		}{
			"正常系: ユーザーが取得できる": {
				email: "test@example.com",
				want:  user,
				err:   nil,
			},
			"異常系: 存在しないメールアドレス": {
				email: "notfound@example.com",
				want:  domain.User{},
				err:   apperr.ErrDataNotFound,
			},
		}

		for name, tt := range tests {
			t.Run(name, func(t *testing.T) {
				got, err := repo.GetByEmail(tt.email)
				if tt.err != nil {
					assert.Error(t, err)
					assert.Equal(t, tt.err, err)
					return
				}

				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			})
		}
	})

	t.Run("Store", func(t *testing.T) {
		defer clearTestData()

		var (
			userID = domain.NewUserID()
			user   = domain.User{
				ID:         userID,
				CognitoUID: "cognito-uid",
				Email:      "test@example.com",
			}
		)

		tests := map[string]struct {
			input *domain.User
			err   error
		}{
			"正常系: ユーザーが作成される": {
				input: &user,
				err:   nil,
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
