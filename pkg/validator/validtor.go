package validator

import "gopkg.in/go-playground/validator.v9"

type Validator struct {
	Validator *validator.Validate
}

func (v *Validator) Validate(i interface{}) error {
	return v.Validator.Struct(i)
}
