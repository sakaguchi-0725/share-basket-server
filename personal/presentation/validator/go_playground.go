package validator

import (
	core "share-basket-server/core/validator"

	"github.com/go-playground/validator"
)

type goPlayground struct {
	v *validator.Validate
}

func (playground *goPlayground) Validate(v any) error {
	return playground.v.Struct(v)
}

func New() core.Validator {
	return &goPlayground{validator.New()}
}
