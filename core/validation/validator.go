package validation

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	goplayground "github.com/go-playground/validator/v10"
)

type Validator interface {
	Validate(interface{}) error
	Messages(err error) []Error
	Register(tag string, fn validator.Func)
	GetTranslator() ut.Translator
	GetValidator() *goplayground.Validate
}

type Error struct {
	Field string      `json:"field"`
	Error interface{} `json:"error"`
}
