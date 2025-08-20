package errors

import "errors"

var (
	ErrInvalidID = errors.New("invalid ID")

	ErrConflict = errors.New("conflict error")

	ErrNotFound = errors.New("not found")

	ErrForeignKeyConstraint = errors.New("cannot delete: foreign key constraint violation")
)
