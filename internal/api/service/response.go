package service

import "errors"

var (
	ErrNotFound     = errors.New("not found")
	ErrBadRequest   = errors.New("bad request")
	ErrInternal     = errors.New("internal errors")
	ErrUnauthorized = errors.New("not authorized")
)
