package repository_test

import (
	"share-basket-server/core/util"
	"share-basket-server/domain"
	"share-basket-server/infra/rdb/repository"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupPersonalShoppingItem() error {
	data := []repository.PersonalShoppingItemDto{
		{ID: 1, Name: "牛乳", Status: domain.UnPurchased.String(), CategoryID: 1, AccountID: dummyAccount.ID},
		{ID: 2, Name: "洗濯洗剤", Status: domain.UnPurchased.String(), CategoryID: 3, AccountID: dummyAccount.ID},
		{ID: 3, Name: "ドッグフード", Status: domain.InTheCart.String(), CategoryID: 4, AccountID: dummyAccount.ID},
		{ID: 4, Name: "卵", Status: domain.InTheCart.String(), CategoryID: 1, AccountID: dummyAccount.ID},
		{ID: 5, Name: "トイレットペーパー", Status: domain.Purchased.String(), CategoryID: 2, AccountID: dummyAccount.ID},
	}

	return testDB.Create(&data).Error
}

func TestPersonalShoppingItemPersistence(t *testing.T) {
	repo := repository.NewPersonalShoppingItemPersistence(testDB)

	t.Run("GetAll", func(t *testing.T) {
		defer clearTestData()

		err := createDummyAccount()
		require.NoError(t, err)

		err = setupPersonalShoppingItem()
		require.NoError(t, err)

		tests := map[string]struct {
			accID  domain.AccountID
			status *domain.ShoppingStatus
			want   []domain.PersonalShoppingItem
			err    error
		}{
			"正常系: ステータスを指定していない場合、全件取得できる": {
				accID:  domain.AccountID(dummyAccount.ID),
				status: nil,
				want: []domain.PersonalShoppingItem{
					{ID: util.Ptr[uint](1), Name: "牛乳", Status: domain.UnPurchased, CategoryID: 1, AccountID: domain.AccountID(dummyAccount.ID)},
					{ID: util.Ptr[uint](2), Name: "洗濯洗剤", Status: domain.UnPurchased, CategoryID: 3, AccountID: domain.AccountID(dummyAccount.ID)},
					{ID: util.Ptr[uint](3), Name: "ドッグフード", Status: domain.InTheCart, CategoryID: 4, AccountID: domain.AccountID(dummyAccount.ID)},
					{ID: util.Ptr[uint](4), Name: "卵", Status: domain.InTheCart, CategoryID: 1, AccountID: domain.AccountID(dummyAccount.ID)},
					{ID: util.Ptr[uint](5), Name: "トイレットペーパー", Status: domain.Purchased, CategoryID: 2, AccountID: domain.AccountID(dummyAccount.ID)},
				},
				err: nil,
			},
			"正常系: 指定したステータスの買い物リストを取得できる": {
				accID:  domain.AccountID(dummyAccount.ID),
				status: util.Ptr(domain.ShoppingStatus(domain.UnPurchased)),
				want: []domain.PersonalShoppingItem{
					{ID: util.Ptr[uint](1), Name: "牛乳", Status: domain.UnPurchased, CategoryID: 1, AccountID: domain.AccountID(dummyAccount.ID)},
					{ID: util.Ptr[uint](2), Name: "洗濯洗剤", Status: domain.UnPurchased, CategoryID: 3, AccountID: domain.AccountID(dummyAccount.ID)},
				},
				err: nil,
			},
			"正常系: 1件も登録されていない場合、空配列が返る": {
				accID:  domain.NewAccountID(),
				status: nil,
				want:   []domain.PersonalShoppingItem{},
				err:    nil,
			},
		}

		for name, tt := range tests {
			t.Run(name, func(t *testing.T) {
				got, err := repo.GetAll(tt.accID, tt.status)

				assert.Equal(t, tt.want, got)
				if tt.err != nil {
					assert.NoError(t, err)
					assert.EqualError(t, tt.err, err.Error())
				} else {
					assert.NoError(t, err)
				}
			})
		}
	})
}
