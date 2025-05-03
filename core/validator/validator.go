package validator

type Validator interface {
	Validate(v any) error
}
