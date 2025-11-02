package request

import (
	"github.com/vault-thirteen/SFRODB/pkg/SFRODB/classes/Method"
	"github.com/vault-thirteen/SFRODB/pkg/SFRODB/classes/UID"
	"github.com/vault-thirteen/SFRODB/pkg/SFRODB/protocol"
)

const (
	ErrSizeIsTooShort = "request size is too short: %v"
	ErrSizeIsTooLong  = "request size is too long: %v"
)

type Request struct {
	Size   uint16 // Length of method name + UID string.
	Method method.Method
	UID    *uid.UID // UID of the requested item.
}

func (r *Request) IsCloseConnection() bool {
	return r.Method == method.Method_CloseConnection
}

func New_CloseConnection() (req *Request, err error) {
	return newSimpleRequest(method.Method_CloseConnection)
}

func New_ResetCache() (req *Request, err error) {
	return newSimpleRequest(method.Method_ResetCache)
}

func New_ShowData(requestedUID string) (req *Request, err error) {
	return newNormalRequest(method.Method_ShowData, requestedUID)
}

func New_SearchRecord(requestedUID string) (req *Request, err error) {
	return newNormalRequest(method.Method_SearchRecord, requestedUID)
}

func New_SearchFile(requestedUID string) (req *Request, err error) {
	return newNormalRequest(method.Method_SearchFile, requestedUID)
}

func New_ForgetRecord(requestedUID string) (req *Request, err error) {
	return newNormalRequest(method.Method_ForgetRecord, requestedUID)
}

func newSimpleRequest(method method.Method) (req *Request, err error) {
	return &Request{
		Size:   protocol.MethodNameLen,
		Method: method,
	}, nil
}

func newNormalRequest(method method.Method, requestedUID string) (req *Request, err error) {
	var u *uid.UID
	u, err = uid.New(requestedUID)
	if err != nil {
		return nil, err
	}

	req = &Request{
		Size:   uint16(protocol.MethodNameLen + u.Length()),
		Method: method,
		UID:    u,
	}

	return req, nil
}
