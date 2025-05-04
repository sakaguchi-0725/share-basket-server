package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"share-basket-server/domain"
)

func TestNewShoppingItem(t *testing.T) {
	categoryID := uint(2)
	status := domain.ShoppingStatus("purchased")

	tests := map[string]struct {
		name       string
		status     *domain.ShoppingStatus
		categoryID uint
		want       domain.ShoppingItem
		err        error
	}{
		"正常系: 商品が正常に作成される": {
			name:       "テスト商品",
			status:     &status,
			categoryID: categoryID,
			want: domain.ShoppingItem{
				Name:       "テスト商品",
				Status:     status,
				CategoryID: categoryID,
			},
			err: nil,
		},
		"正常系: statusがnilの場合、UnPurchasedが設定される": {
			name:       "テスト商品",
			status:     nil,
			categoryID: categoryID,
			want: domain.ShoppingItem{
				Name:       "テスト商品",
				Status:     domain.UnPurchased,
				CategoryID: categoryID,
			},
			err: nil,
		},
		"異常系: 名前が空文字の場合": {
			name:       "",
			status:     &status,
			categoryID: categoryID,
			want:       domain.ShoppingItem{},
			err:        domain.ErrShoppingItemNameRequired,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := domain.NewShoppingItem(tt.name, tt.status, tt.categoryID)

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
