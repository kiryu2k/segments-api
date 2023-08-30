package repository

import "fmt"

var (
	ErrSegmentExists    = fmt.Errorf("specified segment already exists")
	ErrSegmentNotExists = fmt.Errorf("specified segment doesn't exist")
)

var (
	ErrUserExists    = fmt.Errorf("user with specified id already exists")
	ErrUserNotExists = fmt.Errorf("user with specified id doesn't exist")
	ErrHasSegment    = fmt.Errorf("user already has specified segment")
	ErrNoUsers       = fmt.Errorf("there're no users with specified segment")
)
