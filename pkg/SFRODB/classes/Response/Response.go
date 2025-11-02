package response

import (
	"fmt"

	"github.com/vault-thirteen/SFRODB/pkg/SFRODB/classes/Status"
	"github.com/vault-thirteen/SFRODB/pkg/SFRODB/protocol"
)

const (
	ErrContentIsTooLong = "content is too long"
	ErrSizeIsTooShort   = "response size is too short: %v"
	ErrSizeIsTooLong    = "response size is too long: %v"
)

type Response struct {
	Size   uint32 // Length of status name + data array.
	Status status.Status
	Data   []byte
}

func New_ClientError() (resp *Response, err error) {
	return newSimpleResponse(status.Status_ClientError)
}

func New_OK() (resp *Response, err error) {
	return newSimpleResponse(status.Status_OK)
}

func New_ClosingConnection() (resp *Response, err error) {
	return newSimpleResponse(status.Status_ClosingConnection)
}

func New_RecordExists() (resp *Response, err error) {
	return newSimpleResponse(status.Status_RecordExists)
}

func New_RecordDoesNotExist() (resp *Response, err error) {
	return newSimpleResponse(status.Status_RecordDoesNotExist)
}

func New_FileExists() (resp *Response, err error) {
	return newSimpleResponse(status.Status_FileExists)
}

func New_FileDoesNotExist() (resp *Response, err error) {
	return newSimpleResponse(status.Status_FileDoesNotExist)
}

func New_ShowingData(data []byte) (resp *Response, err error) {
	return newNormalResponse(data, status.Status_ShowingData)
}

func newSimpleResponse(status status.Status) (resp *Response, err error) {
	return &Response{
		Size:   protocol.StatusNameLen,
		Status: status,
	}, nil
}

func newNormalResponse(data []byte, status status.Status) (resp *Response, err error) {
	if len(data) > protocol.ContentLenMax {
		return nil, fmt.Errorf(ErrContentIsTooLong)
	}

	resp = &Response{
		Size:   uint32(protocol.StatusNameLen + len(data)),
		Status: status,
		Data:   data,
	}

	return resp, nil
}
