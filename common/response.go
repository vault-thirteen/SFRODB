package common

import (
	"fmt"

	"github.com/vault-thirteen/SFRODB/common/method"
	mn "github.com/vault-thirteen/SFRODB/common/method/name"
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

	// Response Data: Textual and Binary.
	Text string
	Data []byte
}

func newSimpleResponse(method method.Method) (resp *Response, err error) {
	return &Response{
		SRS:           SRS_A,
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

func NewResponse_TextRecordExists() (resp *Response, err error) {
	return newSimpleResponse(method.TextRecordExists)
}

func NewResponse_BinaryRecordExists() (resp *Response, err error) {
	return newSimpleResponse(method.BinaryRecordExists)
}

func NewResponse_TextRecordDoesNotExist() (resp *Response, err error) {
	return newSimpleResponse(method.TextRecordDoesNotExist)
}

func NewResponse_BinaryRecordDoesNotExist() (resp *Response, err error) {
	return newSimpleResponse(method.BinaryRecordDoesNotExist)
}

func NewResponse_TextFileExists() (resp *Response, err error) {
	return newSimpleResponse(method.TextFileExists)
}

func NewResponse_BinaryFileExists() (resp *Response, err error) {
	return newSimpleResponse(method.BinaryFileExists)
}

func NewResponse_TextFileDoesNotExist() (resp *Response, err error) {
	return newSimpleResponse(method.TextFileDoesNotExist)
}

func NewResponse_BinaryFileDoesNotExist() (resp *Response, err error) {
	return newSimpleResponse(method.BinaryFileDoesNotExist)
}

func newNormalResponse(
	text string,
	data []byte,
	useBinary bool,
	method method.Method,
) (resp *Response, err error) {
	resp = &Response{
		SRS:           0, // Will be automatically calculated.
		ResponseSizeC: 0, // Will be automatically calculated.
		Method:        method,
	}

	// Content.
	var contentLen int
	if useBinary {
		resp.Data = data
		contentLen = len(data)
	} else {
		resp.Text = text
		contentLen = len(text)
	}

	// SRS.
	if contentLen > ResponseMessageLengthC-mn.LengthLimit {
		err = fmt.Errorf(ErrTextIsTooLong, ResponseMessageLengthC-mn.LengthLimit, contentLen)
		return nil, err
	} else if contentLen > ResponseMessageLengthB-mn.LengthLimit {
		resp.SRS = SRS_C
	} else if contentLen > ResponseMessageLengthA-mn.LengthLimit {
		resp.SRS = SRS_B
	} else {
		resp.SRS = SRS_A
	}

	// RS.
	switch resp.SRS {
	case SRS_A:
		resp.ResponseSizeA = mn.LengthLimit + uint8(contentLen)
	case SRS_B:
		resp.ResponseSizeB = mn.LengthLimit + uint16(contentLen)
	case SRS_C:
		resp.ResponseSizeC = mn.LengthLimit + uint32(contentLen)
	default:
		return nil, fmt.Errorf(ErrSrsIsNotSupported, resp.SRS)
	}

	return resp, nil
}

func NewResponse_ShowingText(text string) (resp *Response, err error) {
	return newNormalResponse(text, nil, false, method.ShowingText)
}

func NewResponse_ShowingBinary(data []byte) (resp *Response, err error) {
	return newNormalResponse("", data, true, method.ShowingBinary)
}
