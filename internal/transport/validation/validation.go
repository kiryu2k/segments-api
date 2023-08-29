package validation

import (
	"fmt"
	"regexp"
)

const slugMaxSize = 32

var (
	ErrInvalidSize = fmt.Errorf("segment name must be less than %d characters long", slugMaxSize)
	ErrInvalidChar = fmt.Errorf("segment name must consist only word character (alphanumeric & underscore)")
)

func ValidateSlug(slug string) error {
	if len(slug) > slugMaxSize {
		return ErrInvalidSize
	}
	reg, err := regexp.Compile(`^\w+$`)
	if err != nil {
		return err
	}
	if !reg.MatchString(slug) {
		return ErrInvalidChar
	}
	return nil
}
