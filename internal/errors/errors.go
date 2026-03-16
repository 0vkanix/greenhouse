package errors

import "errors"

var (
	ErrInternalServerError = errors.New("the server encountered a problem and could not process your request")
	ErrInvalidIDParam      = errors.New("invalid id parameter")
)
