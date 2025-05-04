package validator

import (
	"github.com/go-playground/validator"
)

type (
	RequestValidator interface {
		Validate(v any) error
	}

	goPlayground struct {
		v *validator.Validate
	}
)

func (playground *goPlayground) Validate(v any) error {
	return playground.v.Struct(v)
}

func New() RequestValidator {
	return &goPlayground{validator.New()}
}
