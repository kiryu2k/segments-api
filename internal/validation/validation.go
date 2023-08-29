package validation

import "fmt"

const slugMaxSize = 32

func ValidateSlug(slug string) error {
	/* TODO: regex validation */
	if len(slug) > slugMaxSize {
		return fmt.Errorf("segment name must be less than %d characters long", slugMaxSize)
	}
	return nil
}
