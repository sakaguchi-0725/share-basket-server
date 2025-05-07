package usecase_test

import (
	"context"
	"errors"
	"testing"

	"share-basket-server/domain"
	. "share-basket-server/test/mock/domain"
	. "share-basket-server/test/mock/usecase"
	"share-basket-server/usecase"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestGetShoppingCategoriesInteractor(t *testing.T) {
	var (
		ctrl   = gomock.NewController(t)
		repo   = NewMockShoppingCategoryRepository(ctrl)
		output = NewMockGetShoppingCategoriesOutputPort(ctrl)
	)

	defer ctrl.Finish()

	tests := map[string]struct {
		setupMock func(repo *MockShoppingCategoryRepository, output *MockGetShoppingCategoriesOutputPort)
		err       error
		want      []usecase.GetShoppingCategoryOutput
	}{
		"正常系: カテゴリの取得に成功": {
			setupMock: func(repo *MockShoppingCategoryRepository, output *MockGetShoppingCategoriesOutputPort) {
				repo.EXPECT().GetAll().Return([]domain.ShoppingCategory{
					{ID: 1, Name: "foods"},
					{ID: 2, Name: "daily"},
				}, nil)
				output.EXPECT().Render(context.Background(), []usecase.GetShoppingCategoryOutput{
					{ID: 1, Name: "foods"},
					{ID: 2, Name: "daily"},
				}).Return(nil)
			},
			err: nil,
			want: []usecase.GetShoppingCategoryOutput{
				{ID: 1, Name: "foods"},
				{ID: 2, Name: "daily"},
			},
		},
		"異常系: リポジトリからエラーが返される": {
			setupMock: func(repo *MockShoppingCategoryRepository, output *MockGetShoppingCategoriesOutputPort) {
				repo.EXPECT().GetAll().Return(nil, errors.New("repository error"))
			},
			err:  errors.New("repository error"),
			want: nil,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			tt.setupMock(repo, output)

			usecase := usecase.NewGetShoppingCategoriesInteractor(repo)
			err := usecase.Execute(context.Background(), output)

			if tt.err == nil {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Equal(t, tt.err.Error(), err.Error())
			}
		})
	}
}
