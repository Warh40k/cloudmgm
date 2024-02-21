package postgres

import "errors"

var (
	ErrUnique   = errors.New("unique costraint violation")
	ErrNoRows   = errors.New("no rows in relation")
	ErrInternal = errors.New("internal error")
)
