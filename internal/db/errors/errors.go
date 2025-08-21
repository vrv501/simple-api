package errors

import "errors"

var (
	ErrInvalidID = errors.New("invalid ID")

	ErrNotFound = errors.New("not found")

	ErrForeignKeyConstraint = errors.New("cannot delete: foreign key constraint violation")
)

type ConflictError struct {
	Key string
	Err error
}

func (e *ConflictError) Error() string {
	return e.Err.Error()
}
