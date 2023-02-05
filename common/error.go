package common

const (
	ErrFileIsNotSet                = "file is not set"
	ErrDataFolderIsNotSet          = "data folder is not set"
	ErrDataFileExtensionIsNotSet   = "data file extension is not set"
	ErrCacheVolumeMaxIsNotSet      = "cache volume limit is not set"
	ErrCachedItemVolumeMaxIsNotSet = "cached item volume limit is not set"
	ErrCachedItemTTLIsNotSet       = "cached item TTL is not set"
	ErrServerHostIsNotSet          = "server host is not set"
	ErrServerPortIsNotSet          = "server port is not set"
	ErrClientHostIsNotSet          = "client host is not set"
	ErrClientPortIsNotSet          = "client port is not set"
	ErrResponseMessageLengthLimit  = "response message length limit is not set"
)

const (
	ErrSrsIsNotSupported = "SRS is unsupported: %d"
	ErrSrsReading        = "SRS reading error: "

	ErrRsReading = "RS reading error: "

	ErrReadingMethodAndData   = "error reading method and data: "
	ErrUnsupportedMethodValue = "unsupported method value: %d"
	ErrUnknownMethodName      = "unknown method name: %s"
	ErrMessageIsTooLong       = "message is too long: %d vs %d"
	ErrTextIsTooLong          = "text is too long: %d vs %d"
	ErrUid                    = "uid error"
	ErrUidIsTooLong           = "uid is too long"

	// ErrSomethingWentWrong is an error which a client sees when he waits for
	// the server's reply and something goes wrong. The reason are variable:
	//	1.	The requested resource has an invalid UID;
	//	2.	The requested resource is not available on the server;
	//	3.	An internal server error has occurred.
	ErrSomethingWentWrong = "something went wrong"
)

type Error struct {
	// Hey you, Go language developers !
	// You wanted Go language to be a replacement for the good old C language.
	// Do you know why C has a "typedef" ? They are smart enough to let all the
	// users use the "type" word in their programs. Shame on you.
	type_  ErrorType
	text   string
	method Method
}

func newError(et ErrorType, msg string, method Method) (err error) {
	return &Error{
		type_:  et,
		text:   msg,
		method: method,
	}
}

func NewServerError(msg string, method Method) (err error) {
	return newError(ErrorTypeServer, msg, method)
}

func NewClientError(msg string, method Method) (err error) {
	return newError(ErrorTypeClient, msg, method)
}

func (e *Error) Error() string {
	return e.text
}

func (e *Error) GetMethod() Method {
	return e.method
}

func (e *Error) IsServerError() bool {
	return e.type_ == ErrorTypeServer
}

func (e *Error) IsClientError() bool {
	return e.type_ == ErrorTypeClient
}
