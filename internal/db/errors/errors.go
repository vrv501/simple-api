package errors

import "errors"

var (
	ErrInvalidID = errors.New("invalid ID")

	ErrNotFound = errors.New("not found")
)

type ConflictError struct {
	Key string
	Err error
}

func (e *ConflictError) Error() string {
	return e.Err.Error()
}

type ForeignKeyError struct {
	Key string
	Err error
}

func (e *ForeignKeyError) Error() string {
	return e.Err.Error()
}
