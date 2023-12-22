package request

import (
	"errors"
	"fmt"

	"github.com/vault-thirteen/SFRODB/pkg/SFRODB/common"
	"github.com/vault-thirteen/SFRODB/pkg/SFRODB/common/error"
	"github.com/vault-thirteen/SFRODB/pkg/SFRODB/common/method"
	"github.com/vault-thirteen/SFRODB/pkg/SFRODB/common/method/name"
	"github.com/vault-thirteen/SFRODB/pkg/SFRODB/common/protocol"
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

func NewRequest_ResetCache() (req *Request, err error) {
	return newSimpleRequest(method.ResetCache)
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

func NewRequest_ShowData(uid string) (req *Request, err error) {
	return newNormalRequest(method.ShowData, uid)
}

func NewRequest_SearchRecord(uid string) (req *Request, err error) {
	return newNormalRequest(method.SearchRecord, uid)
}

func NewRequest_SearchFile(uid string) (req *Request, err error) {
	return newNormalRequest(method.SearchFile, uid)
}

func NewRequest_ForgetRecord(uid string) (req *Request, err error) {
	return newNormalRequest(method.ForgetRecord, uid)
}
