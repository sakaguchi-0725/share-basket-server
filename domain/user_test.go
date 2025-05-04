package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"share-basket-server/domain"
)

func TestNewUser(t *testing.T) {
	tests := map[string]struct {
		id         domain.UserID
		cognitoUID string
		email      string
		want       domain.User
		err        error
	}{
		"正常系: ユーザーが正常に作成される": {
			id:         "test-user-id",
			cognitoUID: "test-cognito-uid",
			email:      "test@example.com",
			want: domain.User{
				ID:         "test-user-id",
				CognitoUID: "test-cognito-uid",
				Email:      "test@example.com",
			},
			err: nil,
		},
		"異常系: CognitoUIDが空文字の場合": {
			id:         "test-user-id",
			cognitoUID: "",
			email:      "test@example.com",
			want:       domain.User{},
			err:        domain.ErrCognitoUIDRequired,
		},
		"異常系: Emailが空文字の場合": {
			id:         "test-user-id",
			cognitoUID: "test-cognito-uid",
			email:      "",
			want:       domain.User{},
			err:        domain.ErrEmailRequired,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := domain.NewUser(tt.id, tt.cognitoUID, tt.email)

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
