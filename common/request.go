package common

import (
	"errors"
	"fmt"
)

type Request struct {
	// Size of the 'Request Size' Field.
	// 'A' => 1 Byte;
	// 'B' => 2 Bytes.
	SRS byte

	// Request Size.
	RequestSizeA uint8
	RequestSizeB uint16

	// Request Method.
	Method Method

	// UID of the Requested Item.
	UID string
}

func (r *Request) IsCloseConnection() bool {
	return r.Method == MethodCloseConnection
}

func newSimpleRequest(method Method) (req *Request, err error) {
	return &Request{
		SRS:          SRS_A,
		RequestSizeA: MethodNameLengthLimit,
		Method:       method,
	}, nil
}

func NewRequest_CloseConnection() (req *Request, err error) {
	return newSimpleRequest(MethodCloseConnection)
}

func NewRequest_ResetTextCache() (req *Request, err error) {
	return newSimpleRequest(MethodResetTextCache)
}

func NewRequest_ResetBinaryCache() (req *Request, err error) {
	return newSimpleRequest(MethodResetBinaryCache)
}

func newNormalRequest(method Method, uid string) (req *Request, err error) {
	if !IsUidValid(uid) {
		return nil, fmt.Errorf(ErrUid)
	}

	req = &Request{
		SRS:          0, // Will be automatically calculated.
		RequestSizeA: 0, // Will be automatically calculated.
		RequestSizeB: 0, // Will be automatically calculated.
		Method:       method,
		UID:          uid,
	}

	// SRS.
	uidLen := len(uid)
	if uidLen <= 0 {
		return nil, errors.New(ErrUid)
	} else if uidLen <= RequestMessageLengthA-MethodNameLengthLimit {
		req.SRS = SRS_A
	} else if uidLen <= RequestMessageLengthB-MethodNameLengthLimit {
		req.SRS = SRS_B
	} else {
		return nil, errors.New(ErrUidIsTooLong)
	}

	// RS.
	switch req.SRS {
	case SRS_A:
		req.RequestSizeA = MethodNameLengthLimit + uint8(uidLen)
	case SRS_B:
		req.RequestSizeB = MethodNameLengthLimit + uint16(uidLen)
	default:
		return nil, fmt.Errorf(ErrSrsIsNotSupported, req.SRS)
	}

	return req, nil
}

func NewRequest_ShowText(uid string) (req *Request, err error) {
	return newNormalRequest(MethodShowText, uid)
}

func NewRequest_ShowBinary(uid string) (req *Request, err error) {
	return newNormalRequest(MethodShowBinary, uid)
}

func NewRequest_SearchTextRecord(uid string) (req *Request, err error) {
	return newNormalRequest(MethodSearchTextRecord, uid)
}

func NewRequest_SearchBinaryRecord(uid string) (req *Request, err error) {
	return newNormalRequest(MethodSearchBinaryRecord, uid)
}

func NewRequest_SearchTextFile(uid string) (req *Request, err error) {
	return newNormalRequest(MethodSearchTextFile, uid)
}

func NewRequest_SearchBinaryFile(uid string) (req *Request, err error) {
	return newNormalRequest(MethodSearchBinaryFile, uid)
}

func NewRequest_ForgetTextRecord(uid string) (req *Request, err error) {
	return newNormalRequest(MethodForgetTextRecord, uid)
}

func NewRequest_ForgetBinaryRecord(uid string) (req *Request, err error) {
	return newNormalRequest(MethodForgetBinaryRecord, uid)
}
