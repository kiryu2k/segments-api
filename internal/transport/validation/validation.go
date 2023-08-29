package validation

import (
	"fmt"
	"regexp"

	"github.com/kiryu-dev/segments-api/pkg/util/parser"
)

const slugMaxSize = 32

var (
	ErrInvalidSize = fmt.Errorf("segment name must be less than %d characters long", slugMaxSize)
	ErrInvalidChar = fmt.Errorf("segment name must consist only word character (alphanumeric & underscore)")
	ErrRegexpErr   = fmt.Errorf("unexpected regexp error")
)

func ValidateSlug(slug string) error {
	if len(slug) > slugMaxSize {
		return ErrInvalidSize
	}
	match, err := regexp.MatchString(`^\w+$`, slug)
	if err != nil {
		return err
	}
	if !match {
		return ErrInvalidChar
	}
	return nil
}

func ValidateTTL(ttl string) (*parser.TTL, error) {
	patterns := [...]string{
		`^(\d+y)(\d+m)?(\d+d)?$`,
		`^(\d+y)?(\d+m)(\d+d)?$`,
		`^(\d+y)?(\d+m)?(\d+d)$`,
	}
	for _, pattern := range patterns {
		reg, err := regexp.Compile(pattern)
		if err != nil {
			return nil, ErrRegexpErr
		}
		if reg.MatchString(ttl) {
			return parser.ParseTTL(ttl)
		}
	}
	return nil, fmt.Errorf("invalid ttl format: expected something like this 1y8m16d")
}
