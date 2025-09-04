package errors

import "errors"

var (
	ErrInvalidValue = errors.New("invalid value")

	ErrNotFound = errors.New("not found")

	ErrConflict = errors.New("conflict")

	ErrForeignKeyViolation = errors.New("foreign key constraint failed")
)

type HintError struct {
	Key string
	Err error
}

func (e *HintError) Error() string {
	return e.Err.Error()
}
