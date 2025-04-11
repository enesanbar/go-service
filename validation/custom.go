package validation

import "github.com/go-playground/validator/v10"

type CustomValidation struct {
	Tag      string
	Func     validator.Func
	Messages map[string]string
}
