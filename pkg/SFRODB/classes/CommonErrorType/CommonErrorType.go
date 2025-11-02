package cet

// CommonErrorType is a common error type.
type CommonErrorType byte

const (
	CommonErrorType_Server = CommonErrorType(1)
	CommonErrorType_Client = CommonErrorType(2)
)
