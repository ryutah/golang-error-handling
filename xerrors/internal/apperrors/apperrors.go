package apperrors

import "errors"

var (
	ErrNotFound         = errors.New("not found")
	ErrInvalidParameter = errors.New("invalid parameter")
	ErrInternal         = errors.New("internal")
)
