package database_test

import (
	"share-basket-server/personal/domain"
	"share-basket-server/personal/infra/database"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShoppingCategoryPersistence(t *testing.T) {
	repo := database.NewShoppingCategoryPersistence(testDB)

	t.Run("GetAll", func(t *testing.T) {
		defer clearTestData()

		tests := map[string]struct {
			want []domain.ShoppingCategory
			err  error
		}{
			"正常系: カテゴリーが取得できる": {
				want: []domain.ShoppingCategory{
					domain.NewShoppingCategory(1, "foods"),
					domain.NewShoppingCategory(2, "daily"),
					domain.NewShoppingCategory(3, "hygiene"),
					domain.NewShoppingCategory(4, "pet"),
					domain.NewShoppingCategory(5, "healthcare"),
					domain.NewShoppingCategory(6, "miscellaneous"),
					domain.NewShoppingCategory(7, "hobby"),
				},
				err: nil,
			},
		}

		for name, tt := range tests {
			t.Run(name, func(t *testing.T) {
				got, err := repo.GetAll()
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
}
