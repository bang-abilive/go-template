package values

import (
	"errors"
	"strings"
)

type Email struct {
	value string
}

func NewEmail(value string) (Email, error) {
	value = strings.TrimSpace(value)
	if !strings.Contains(value, "@") {
		return Email{}, errors.New("invalid email format")
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
