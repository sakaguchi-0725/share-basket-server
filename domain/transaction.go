//go:generate mockgen -destination=../mock/$GOPACKAGE/$GOFILE . Transaction
package domain

import "context"

type Transaction interface {
	Run(ctx context.Context, fn func(ctx context.Context) error) error
}
