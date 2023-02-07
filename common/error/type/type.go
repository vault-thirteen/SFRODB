package cet

// ErrorType is a common error type.
type ErrorType byte

const (
	ErrorTypeServer = ErrorType(1)
	ErrorTypeClient = ErrorType(2)
)
