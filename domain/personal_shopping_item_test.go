package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"share-basket-server/core/util"
	"share-basket-server/domain"
)

func TestNewPersonalShoppingItem(t *testing.T) {
	categoryID := uint(2)
	accID := domain.NewAccountID()

	tests := map[string]struct {
		name       string
		status     *domain.ShoppingStatus
		accID      domain.AccountID
		categoryID uint
		want       domain.PersonalShoppingItem
		err        error
	}{
		"正常系: 商品が正常に作成される": {
			name:       "テスト商品",
			status:     util.Ptr(domain.UnPurchased),
			categoryID: categoryID,
			accID:      accID,
			want: domain.PersonalShoppingItem{
				Name:       "テスト商品",
				Status:     domain.UnPurchased,
				CategoryID: categoryID,
				AccountID:  accID,
			},
			err: nil,
		},
		"正常系: statusがnilの場合、UnPurchasedが設定される": {
			name:       "テスト商品",
			status:     nil,
			categoryID: categoryID,
			accID:      accID,
			want: domain.PersonalShoppingItem{
				Name:       "テスト商品",
				Status:     domain.UnPurchased,
				CategoryID: categoryID,
				AccountID:  accID,
			},
			err: nil,
		},
		"異常系: 名前が空文字の場合": {
			name:       "",
			status:     util.Ptr(domain.InTheCart),
			categoryID: categoryID,
			accID:      accID,
			want:       domain.PersonalShoppingItem{},
			err:        domain.ErrPersonalShoppingItemNameRequired,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := domain.NewPersonalShoppingItem(tt.name, tt.status, tt.categoryID, tt.accID)

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
