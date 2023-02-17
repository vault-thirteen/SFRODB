package ce

import (
	"github.com/vault-thirteen/SFRODB/pkg/common/error/type"
	"github.com/vault-thirteen/SFRODB/pkg/common/method"
)

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
	ErrSrsIsNotSupported      = "SRS is unsupported: %d"
	ErrSrsReading             = "SRS reading error: "
	ErrRsReading              = "RS reading error: "
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

// CommonError is a common error used by the service.
type CommonError struct {
	// Hey you, Go language developers !
	// You wanted Go language to be a replacement for the good old C language.
	// Do you know why C has a "typedef" ? They are smart enough to let all the
	// users use the "type" word in their programs. Shame on you.
	typeOfError cet.ErrorType
	text        string
	method      method.Method

	// ID of a client that created this error.
	clientId string
}

func newCommonError(
	et cet.ErrorType,
	msg string,
	method method.Method,
	clientId string,
) (ce *CommonError) {
	return &CommonError{
		typeOfError: et,
		text:        msg,
		method:      method,
		clientId:    clientId,
	}
}

func NewServerError(
	msg string,
	method method.Method,
	clientId string,
) (ce *CommonError) {
	return newCommonError(cet.ErrorTypeServer, msg, method, clientId)
}

func NewClientError(
	msg string,
	method method.Method,
	clientId string,
) (ce *CommonError) {
	return newCommonError(cet.ErrorTypeClient, msg, method, clientId)
}

func (ce *CommonError) Error() string {
	return ce.text
}

func (ce *CommonError) GetMethod() method.Method {
	return ce.method
}

func (ce *CommonError) IsServerError() bool {
	return ce.typeOfError == cet.ErrorTypeServer
}

func (ce *CommonError) IsClientError() bool {
	return ce.typeOfError == cet.ErrorTypeClient
}

func (ce *CommonError) GetClientId() (clientId string) {
	return ce.clientId
}
