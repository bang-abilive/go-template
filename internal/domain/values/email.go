package values

import (
	"errors"
	"fmt"
	"strings"
)

type Email struct {
	value string
}

var ErrInvalidEmailFormat = errors.New("[values] email invalid")

func NewEmail(value string) (Email, error) {
	value = strings.TrimSpace(value)
	if !strings.Contains(value, "@") {
		return Email{}, fmt.Errorf("%w: %s", ErrInvalidEmailFormat, value)
	}
	// Trả về một instance mới của Email
	return Email{value: value}, nil
}

func (e Email) String() string {
	return e.value
}

func (e Email) Value() string {
	return e.value
}

func (e Email) IsValid() bool {
	return strings.Contains(e.value, "@")
}
