package output

import "context"

type LoginOutputPort interface {
	Render(ctx context.Context, token string) error
}
