package common

import "github.com/go-playground/validator/v10"

type EchoValidator struct {
	Validator *validator.Validate
}

func (v *EchoValidator) Validate(i interface{}) error {
	return v.Validator.Struct(i)
}
