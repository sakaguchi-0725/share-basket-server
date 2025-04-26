package output

import "context"

type SignUpOutputPort interface {
	Render(ctx context.Context) error
}
