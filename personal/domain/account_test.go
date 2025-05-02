package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"share-basket-server/personal/domain"
)

func TestNewAccount(t *testing.T) {
	tests := map[string]struct {
		id     domain.AccountID
		userID domain.UserID
		name   string
		want   domain.Account
		err    error
	}{
		"正常系: アカウントが正常に作成される": {
			id:     "test-account-id",
			userID: "test-user-id",
			name:   "テストアカウント",
			want: domain.Account{
				ID:     "test-account-id",
				UserID: "test-user-id",
				Name:   "テストアカウント",
			},
			err: nil,
		},
		"異常系: 名前が空文字の場合": {
			id:     "test-account-id",
			userID: "test-user-id",
			name:   "",
			want:   domain.Account{},
			err:    domain.ErrAccountNameRequired,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := domain.NewAccount(tt.id, tt.userID, tt.name)

			if tt.err != nil {
				require.Error(t, err)
				assert.EqualError(t, tt.err, err.Error())
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
