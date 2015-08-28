package ql

import "errors"

// package errors
var (
	ErrTableNotSpecified  = errors.New("ql: table not specified")
	ErrColumnNotSpecified = errors.New("ql: column not specified")
	ErrLoadNonPointer     = errors.New("ql: attempt to load into a non-pointer")
	ErrPlaceholderCount   = errors.New("ql: wrong placeholder count")
	ErrNotSupported       = errors.New("ql: not supported")
)
