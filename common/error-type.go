package common

type ErrorType byte

const (
	ErrorTypeServer = ErrorType(1)
	ErrorTypeClient = ErrorType(2)
)
