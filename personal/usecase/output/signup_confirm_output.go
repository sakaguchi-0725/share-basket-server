package output

import "context"

type SignUpConfirmOutputPort interface {
	Render(ctx context.Context) error
}
