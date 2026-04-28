package validation

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
)

type Validation struct {
	validator *validator.Validate
}

func NewValidation() *Validation {
	return &Validation{
		validator: validator.New(validator.WithRequiredStructEnabled()),
	}
}

func (v *Validation) Validate(s any) error {
	if err := v.validator.Struct(s); err != nil {
		return echo.ErrBadRequest.Wrap(err)
	}
	return nil
}
