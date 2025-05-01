package input

import (
	"context"
	"share-basket-server/personal/usecase/output"
)

type GetShoppingCategoriesPort interface {
	Execute(ctx context.Context, out output.GetShoppingCategoriesPort) error
}
