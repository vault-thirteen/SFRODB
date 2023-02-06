package common

import (
	"fmt"
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
	Method Method

	// Response Data: Textual and Binary.
	Text string
	Data []byte
}

func newSimpleResponse(method Method) (resp *Response, err error) {
	return &Response{
		SRS:           SRS_A,
		ResponseSizeA: MethodNameLengthLimit,
		Method:        method,
	}, nil
}

func NewResponse_ClientError() (resp *Response, err error) {
	return newSimpleResponse(MethodClientError)
}

func NewResponse_OK() (resp *Response, err error) {
	return newSimpleResponse(MethodOK)
}

func NewResponse_ClosingConnection() (resp *Response, err error) {
	return newSimpleResponse(MethodClosingConnection)
}

func NewResponse_TextRecordExists() (resp *Response, err error) {
	return newSimpleResponse(MethodTextRecordExists)
}

func NewResponse_BinaryRecordExists() (resp *Response, err error) {
	return newSimpleResponse(MethodBinaryRecordExists)
}

func NewResponse_TextRecordDoesNotExist() (resp *Response, err error) {
	return newSimpleResponse(MethodTextRecordDoesNotExist)
}

func NewResponse_BinaryRecordDoesNotExist() (resp *Response, err error) {
	return newSimpleResponse(MethodBinaryRecordDoesNotExist)
}

func NewResponse_TextFileExists() (resp *Response, err error) {
	return newSimpleResponse(MethodTextFileExists)
}

func NewResponse_BinaryFileExists() (resp *Response, err error) {
	return newSimpleResponse(MethodBinaryFileExists)
}

func NewResponse_TextFileDoesNotExist() (resp *Response, err error) {
	return newSimpleResponse(MethodTextFileDoesNotExist)
}

func NewResponse_BinaryFileDoesNotExist() (resp *Response, err error) {
	return newSimpleResponse(MethodBinaryFileDoesNotExist)
}

func newNormalResponse(
	text string,
	data []byte,
	useBinary bool,
	method Method,
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
	if contentLen > ResponseMessageLengthC-MethodNameLengthLimit {
		err = fmt.Errorf(ErrTextIsTooLong, ResponseMessageLengthC-MethodNameLengthLimit, contentLen)
		return nil, err
	} else if contentLen > ResponseMessageLengthB-MethodNameLengthLimit {
		resp.SRS = SRS_C
	} else if contentLen > ResponseMessageLengthA-MethodNameLengthLimit {
		resp.SRS = SRS_B
	} else {
		resp.SRS = SRS_A
	}

	// RS.
	switch resp.SRS {
	case SRS_A:
		resp.ResponseSizeA = MethodNameLengthLimit + uint8(contentLen)
	case SRS_B:
		resp.ResponseSizeB = MethodNameLengthLimit + uint16(contentLen)
	case SRS_C:
		resp.ResponseSizeC = MethodNameLengthLimit + uint32(contentLen)
	default:
		return nil, fmt.Errorf(ErrSrsIsNotSupported, resp.SRS)
	}

	return resp, nil
}

func NewResponse_ShowingText(text string) (resp *Response, err error) {
	return newNormalResponse(text, nil, false, MethodShowingText)
}

func NewResponse_ShowingBinary(data []byte) (resp *Response, err error) {
	return newNormalResponse("", data, true, MethodShowingBinary)
}
