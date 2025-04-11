package config

import (
	"fmt"
)

const (
	TemplateMissingProperty = "missing property: '%s'"
)

type ErrMissingProperty struct {
	err string
}

func NewErrMissingProperty(err string) *ErrMissingProperty {
	return &ErrMissingProperty{err: err}
}

func (e ErrMissingProperty) Error() string {
	return fmt.Sprintf(TemplateMissingProperty, e.err)
}

func NewMissingPropertyError(property string) error {
	return fmt.Errorf(TemplateMissingProperty, property)
}
