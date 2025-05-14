package repository

import (
	"context"
)

type Transaction interface {
	WithTx(ctx context.Context, fn func(ctx context.Context) error) error
}
