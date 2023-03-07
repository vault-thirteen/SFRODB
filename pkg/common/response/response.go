package response

import (
	"fmt"

	"github.com/vault-thirteen/SFRODB/pkg/common/error"
	"github.com/vault-thirteen/SFRODB/pkg/common/method"
	"github.com/vault-thirteen/SFRODB/pkg/common/method/name"
	"github.com/vault-thirteen/SFRODB/pkg/common/protocol"
)

type Response struct {
	// Size of the 'Response Size' Field.
	// 'C' => 4 Bytes.
	SRS byte

	// Response Size.
	ResponseSizeA byte
	ResponseSizeB uint16
	ResponseSizeC uint32

	// Response Method.
	Method method.Method

	// Response Data.
	Data []byte
}

func newSimpleResponse(method method.Method) (resp *Response, err error) {
	return &Response{
		SRS:           proto.SRS_A,
		ResponseSizeA: mn.LengthLimit,
		Method:        method,
	}, nil
}

func NewResponse_ClientError() (resp *Response, err error) {
	return newSimpleResponse(method.ClientError)
}

func NewResponse_OK() (resp *Response, err error) {
	return newSimpleResponse(method.OK)
}

func NewResponse_ClosingConnection() (resp *Response, err error) {
	return newSimpleResponse(method.ClosingConnection)
}

func NewResponse_RecordExists() (resp *Response, err error) {
	return newSimpleResponse(method.RecordExists)
}

func NewResponse_RecordDoesNotExist() (resp *Response, err error) {
	return newSimpleResponse(method.RecordDoesNotExist)
}

func NewResponse_FileExists() (resp *Response, err error) {
	return newSimpleResponse(method.FileExists)
}

func NewResponse_FileDoesNotExist() (resp *Response, err error) {
	return newSimpleResponse(method.FileDoesNotExist)
}

func newNormalResponse(
	data []byte,
	method method.Method,
) (resp *Response, err error) {
	resp = &Response{
		SRS:           0, // Will be automatically calculated.
		ResponseSizeC: 0, // Will be automatically calculated.
		Method:        method,
		Data:          data,
	}

	var contentLen = len(data)

	// SRS.
	if contentLen > proto.ResponseMessageLengthC-mn.LengthLimit {
		err = fmt.Errorf(ce.ErrTextIsTooLong, proto.ResponseMessageLengthC-mn.LengthLimit, contentLen)
		return nil, err
	} else if contentLen > proto.ResponseMessageLengthB-mn.LengthLimit {
		resp.SRS = proto.SRS_C
	} else if contentLen > proto.ResponseMessageLengthA-mn.LengthLimit {
		resp.SRS = proto.SRS_B
	} else {
		resp.SRS = proto.SRS_A
	}

	// RS.
	switch resp.SRS {
	case proto.SRS_A:
		resp.ResponseSizeA = mn.LengthLimit + uint8(contentLen)
	case proto.SRS_B:
		resp.ResponseSizeB = mn.LengthLimit + uint16(contentLen)
	case proto.SRS_C:
		resp.ResponseSizeC = mn.LengthLimit + uint32(contentLen)
	default:
		return nil, fmt.Errorf(ce.ErrSrsIsNotSupported, resp.SRS)
	}

	return resp, nil
}

func NewResponse_ShowingData(data []byte) (resp *Response, err error) {
	return newNormalResponse(data, method.ShowingData)
}
