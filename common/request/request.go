package request

import (
	"errors"
	"fmt"

	"github.com/vault-thirteen/SFRODB/common"
	ce "github.com/vault-thirteen/SFRODB/common/error"
	"github.com/vault-thirteen/SFRODB/common/method"
	"github.com/vault-thirteen/SFRODB/common/method/name"
	"github.com/vault-thirteen/SFRODB/common/protocol"
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
	Method method.Method

	// UID of the Requested Item.
	UID string
}

func (r *Request) IsCloseConnection() bool {
	return r.Method == method.CloseConnection
}

func newSimpleRequest(method method.Method) (req *Request, err error) {
	return &Request{
		SRS:          proto.SRS_A,
		RequestSizeA: mn.LengthLimit,
		Method:       method,
	}, nil
}

func NewRequest_CloseConnection() (req *Request, err error) {
	return newSimpleRequest(method.CloseConnection)
}

func NewRequest_ResetTextCache() (req *Request, err error) {
	return newSimpleRequest(method.ResetTextCache)
}

func NewRequest_ResetBinaryCache() (req *Request, err error) {
	return newSimpleRequest(method.ResetBinaryCache)
}

func newNormalRequest(method method.Method, uid string) (req *Request, err error) {
	if !common.IsUidValid(uid) {
		return nil, fmt.Errorf(ce.ErrUid)
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
		return nil, errors.New(ce.ErrUid)
	} else if uidLen <= proto.RequestMessageLengthA-mn.LengthLimit {
		req.SRS = proto.SRS_A
	} else if uidLen <= proto.RequestMessageLengthB-mn.LengthLimit {
		req.SRS = proto.SRS_B
	} else {
		return nil, errors.New(ce.ErrUidIsTooLong)
	}

	// RS.
	switch req.SRS {
	case proto.SRS_A:
		req.RequestSizeA = mn.LengthLimit + uint8(uidLen)
	case proto.SRS_B:
		req.RequestSizeB = mn.LengthLimit + uint16(uidLen)
	default:
		return nil, fmt.Errorf(ce.ErrSrsIsNotSupported, req.SRS)
	}

	return req, nil
}

func NewRequest_ShowText(uid string) (req *Request, err error) {
	return newNormalRequest(method.ShowText, uid)
}

func NewRequest_ShowBinary(uid string) (req *Request, err error) {
	return newNormalRequest(method.ShowBinary, uid)
}

func NewRequest_SearchTextRecord(uid string) (req *Request, err error) {
	return newNormalRequest(method.SearchTextRecord, uid)
}

func NewRequest_SearchBinaryRecord(uid string) (req *Request, err error) {
	return newNormalRequest(method.SearchBinaryRecord, uid)
}

func NewRequest_SearchTextFile(uid string) (req *Request, err error) {
	return newNormalRequest(method.SearchTextFile, uid)
}

func NewRequest_SearchBinaryFile(uid string) (req *Request, err error) {
	return newNormalRequest(method.SearchBinaryFile, uid)
}

func NewRequest_ForgetTextRecord(uid string) (req *Request, err error) {
	return newNormalRequest(method.ForgetTextRecord, uid)
}

func NewRequest_ForgetBinaryRecord(uid string) (req *Request, err error) {
	return newNormalRequest(method.ForgetBinaryRecord, uid)
}
