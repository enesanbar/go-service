package validation

import (
	"errors"
	"strings"

	"go.uber.org/fx"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"

	goplayground "github.com/go-playground/validator/v10"
	translations "github.com/go-playground/validator/v10/translations/en"
)

type goPlayground struct {
	validator *goplayground.Validate
	Translate ut.Translator
}

func (g *goPlayground) GetValidator() *goplayground.Validate {
	return g.validator
}

func (g *goPlayground) GetTranslator() ut.Translator {
	return g.Translate
}

type Params struct {
	fx.In

	CustomValidators []CustomValidation `group:"validators"`
}

// NewGoPlayground returns go-playground implementation of Validator interface
func NewGoPlayground(p Params) (Validator, error) {
	var (
		language         = en.New()
		uni              = ut.New(language, language)
		translate, found = uni.GetTranslator("en")
	)

	if !found {
		return nil, errors.New("translator not found")
	}

	v := goplayground.New()
	if err := translations.RegisterDefaultTranslations(v, translate); err != nil {
		return nil, errors.New("translator not found")
	}

	for _, cv := range p.CustomValidators {
		v.RegisterValidation(cv.Tag, cv.Func)
		v.RegisterTranslation(cv.Tag, translate, func(ut ut.Translator) error {
			return ut.Add(cv.Tag, cv.Messages["en"], true)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T(cv.Tag, fe.Field())

			return t
		})
	}

	return &goPlayground{validator: v, Translate: translate}, nil
}

func (g *goPlayground) Validate(i interface{}) error {
	return g.validator.Struct(i)
}

func (g *goPlayground) Messages(rawError error) []Error {
	errs := make([]Error, 0)
	validationErrors := rawError.(goplayground.ValidationErrors)
	for _, validationError := range validationErrors {
		errs = append(errs, Error{
			strings.ToLower(validationError.Field()),
			validationError.Translate(g.Translate),
		})
	}

	return errs
}

func (g *goPlayground) Register(tag string, fn goplayground.Func) {
	g.validator.RegisterValidation(tag, fn)
}
