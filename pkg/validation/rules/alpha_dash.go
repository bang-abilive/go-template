package rules

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

// Example: manage.user_create, auth.login, role.read_all
var alphaDashRegex = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]*(\.[a-zA-Z][a-zA-Z0-9_]*)*$`)

func AlphaDash(fl validator.FieldLevel) bool {
	return alphaDashRegex.MatchString(fl.Field().String())
}
