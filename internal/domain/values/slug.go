package values

import "github.com/samber/lo"

type Slug struct {
	value string
}

func NewSlug(value string) (Slug, error) {
	return Slug{value: lo.KebabCase(value)}, nil
}

func (s Slug) String() string {
	return s.value
}

func (s Slug) Value() string {
	return s.value
}

func (s Slug) IsValid() bool {
	return len(s.value) >= 3 && len(s.value) <= 64
}
