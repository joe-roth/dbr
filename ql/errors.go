package ql

import "errors"

// package errors
var (
	ErrTableNotSpecified  = errors.New("ql: table not specified")
	ErrColumnNotSpecified = errors.New("ql: column not specified")
	ErrBadArgument        = errors.New("ql: bad argument")
	ErrNotSupported       = errors.New("ql: not supported")
)
