package ce

import (
	"github.com/vault-thirteen/SFRODB/pkg/SFRODB/classes/CommonErrorType"
	"github.com/vault-thirteen/SFRODB/pkg/SFRODB/classes/Method"
	status "github.com/vault-thirteen/SFRODB/pkg/SFRODB/classes/Status"
)

// CommonError is a common error used by the service.
type CommonError struct {
	typé     cet.CommonErrorType
	text     string
	method   method.Method
	status   status.Status
	clientId string // ID of a client that created this error.
}

func newCommonError(
	typé cet.CommonErrorType,
	msg string,
	method method.Method,
	status status.Status,
	clientId string,
) (ce *CommonError) {
	return &CommonError{
		typé:     typé,
		text:     msg,
		method:   method,
		status:   status,
		clientId: clientId,
	}
}

func NewServerError(
	msg string,
	method method.Method,
	status status.Status,
	clientId string,
) (ce *CommonError) {
	return newCommonError(cet.CommonErrorType_Server, msg, method, status, clientId)
}

func NewClientError(
	msg string,
	method method.Method,
	status status.Status,
	clientId string,
) (ce *CommonError) {
	return newCommonError(cet.CommonErrorType_Client, msg, method, status, clientId)
}

func (ce *CommonError) Error() string {
	return ce.text
}

func (ce *CommonError) GetMethod() method.Method {
	return ce.method
}

func (ce *CommonError) IsServerError() bool {
	return ce.typé == cet.CommonErrorType_Server
}

func (ce *CommonError) IsClientError() bool {
	return ce.typé == cet.CommonErrorType_Client
}

func (ce *CommonError) GetClientId() (clientId string) {
	return ce.clientId
}
