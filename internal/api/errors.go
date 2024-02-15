package api

import "errors"

var (
	ErrNotFound   = errors.New("not found")
	ErrBadRequest = errors.New("bad request")
	ErrInternal   = errors.New("internal errors")
)

type NotFoundError struct {
	Message string `json:"message"`
	Err     error
}

func (e NotFoundError) Error() string {
	return e.Message
}

func (e NotFoundError) Unwrap() error {
	return e.Err
}

type BadRequestError struct {
	Message string `json:"message"`
	Err     error
}

func (e BadRequestError) Error() string {
	return e.Message
}

func (e BadRequestError) Unwrap() error {
	return e.Err
}

type InternalError struct {
	Message string `json:"message"`
	Err     error
}

func (e InternalError) Error() string {
	return e.Message
}

func (e InternalError) Unwrap() error {
	return e.Err
}
