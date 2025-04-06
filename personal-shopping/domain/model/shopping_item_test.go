package model_test

import (
	"errors"
	"share-basket/personal-shopping/core/util"
	"share-basket/personal-shopping/domain/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShoppingItem(t *testing.T) {
	t.Run("NewShoppingItem", func(t *testing.T) {
		tests := map[string]struct {
			name        string
			category    model.ShoppingCategory
			status      model.ShoppingStatus
			expected    model.ShoppingItem
			expectedErr error
		}{
			"正常系": {
				name:     "牛乳",
				category: model.NewShoppingCategory(util.Ptr[int64](1), "foods"),
				status:   model.UnPurchased,
				expected: model.ShoppingItem{
					ID:       nil,
					Name:     "牛乳",
					Category: model.NewShoppingCategory(util.Ptr[int64](1), "foods"),
					Status:   model.UnPurchased,
				},
				expectedErr: nil,
			},
			"買い物リスト名が空文字の場合": {
				name:        "",
				category:    model.NewShoppingCategory(util.Ptr[int64](1), "foods"),
				status:      model.UnPurchased,
				expected:    model.ShoppingItem{},
				expectedErr: errors.New("shopping name is required"),
			},
		}

		for name, tt := range tests {
			t.Run(name, func(t *testing.T) {
				actual, err := model.NewShoppingItem(tt.name, tt.category, tt.status)

				assert.Equal(t, tt.expected, actual)
				if tt.expectedErr != nil {
					assert.Error(t, err)
					assert.EqualError(t, err, tt.expectedErr.Error())
				} else {
					assert.NoError(t, err)
				}
			})
		}
	})

	t.Run("RecreateShoppingItem", func(t *testing.T) {
		tests := map[string]struct {
			id       *int64
			name     string
			category model.ShoppingCategory
			status   model.ShoppingStatus
			expected model.ShoppingItem
		}{
			"正常系": {
				id:       util.Ptr[int64](1),
				name:     "牛乳",
				category: model.NewShoppingCategory(util.Ptr[int64](1), "foods"),
				status:   model.UnPurchased,
				expected: model.ShoppingItem{
					ID:       util.Ptr[int64](1),
					Name:     "牛乳",
					Category: model.NewShoppingCategory(util.Ptr[int64](1), "foods"),
					Status:   model.UnPurchased,
				},
			},
		}

		for name, tt := range tests {
			t.Run(name, func(t *testing.T) {
				actual := model.RecreateShoppingItem(tt.id, tt.name, tt.category, tt.status)
				assert.Equal(t, tt.expected, actual)
			})
		}
	})
}
