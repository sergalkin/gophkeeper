package apperr

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidInput         = errors.New("invalid input")
	ErrConflict             = fmt.Errorf("conflict: %w", ErrInvalidInput)
	ErrUpdatedAtDoesntMatch = fmt.Errorf("could not update secrete. Local data doesn't match with server")
)
