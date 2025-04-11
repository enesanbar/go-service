package validation

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	go_playground "github.com/go-playground/validator/v10"
)

type Validator interface {
	Validate(interface{}) error
	Messages(err error) []Error
	Register(tag string, fn validator.Func)
	GetTranslator() ut.Translator
	GetValidator() *go_playground.Validate
}

type Error struct {
	Field string      `json:"field"`
	Error interface{} `json:"error"`
}
