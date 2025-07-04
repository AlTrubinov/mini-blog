package apperror

import (
	"errors"
	"net/http"
)

var (
	ErrValidation   = errors.New("validation error")
	ErrNotFound     = errors.New("not found")
	ErrForbidden    = errors.New("forbidden")
	ErrInternal     = errors.New("internal error")
	ErrTimeout      = errors.New("timeout error")
	ErrUnauthorized = errors.New("authorization failed")
)

func GetCodeByError(err error) int {
	var code int
	switch {
	case errors.Is(err, ErrValidation):
		code = http.StatusBadRequest
	case errors.Is(err, ErrNotFound):
		code = http.StatusNotFound
	case errors.Is(err, ErrForbidden):
		code = http.StatusForbidden
	case errors.Is(err, ErrTimeout):
		code = http.StatusRequestTimeout
	case errors.Is(err, ErrUnauthorized):
		code = http.StatusUnauthorized
	default:
		code = http.StatusInternalServerError
	}

	return code
}
