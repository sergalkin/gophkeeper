package apperr

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidInput = errors.New("invalid input")
	ErrConflict     = fmt.Errorf("conflict: %w", ErrInvalidInput)
	ErrNotFound     = fmt.Errorf("not found: %w", ErrInvalidInput)
)
