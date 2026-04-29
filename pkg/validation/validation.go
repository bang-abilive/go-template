package validation

import (
	"ndinhbang/go-template/pkg/validation/rules"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
)

type Validation struct {
	validator *validator.Validate
}

func NewValidation() *Validation {
	v := validator.New(validator.WithRequiredStructEnabled())
	// Register custom validation rules
	v.RegisterValidation("alpha_dash", rules.AlphaDash)
	return &Validation{
		validator: v,
	}
}

func (v *Validation) Validate(s any) error {
	if err := v.validator.Struct(s); err != nil {
		return echo.ErrBadRequest.Wrap(err)
	}
	return nil
}
