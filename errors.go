package dbr

import "errors"

// package errors
var (
	ErrNotFound     = errors.New("dbr: not found")
	ErrNotSupported = errors.New("dbr: not supported")
)
