package errors

import "errors"

var (
	ErrInvalidId = errors.New("invalid ID")

	ErrConflict = errors.New("conflict error")

	ErrNotFound = errors.New("not found")
)
