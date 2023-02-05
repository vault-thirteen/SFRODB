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

func (r *Request) IsShowText() bool {
	return r.Method == MethodShowText
}

func (r *Request) IsShowBinary() bool {
	return r.Method == MethodShowBinary
}

func NewRequest_CloseConnection() (req *Request, err error) {
	return &Request{
		SRS:          SRS_A,
		RequestSizeA: MethodNameLengthLimit,
		Method:       MethodCloseConnection,
	}, nil
}

func NewRequest_ShowText(uid string) (req *Request, err error) {
	if !IsUidValid(uid) {
		return nil, fmt.Errorf(ErrUid)
	}

	req = &Request{
		SRS:          0, // Will be automatically calculated.
		RequestSizeA: 0, // Will be automatically calculated.
		RequestSizeB: 0, // Will be automatically calculated.
		Method:       MethodShowText,
		UID:          uid,
	}

	err = req.calculateSRS(uid)
	if err != nil {
		return nil, err
	}

	err = req.calculateRequestSize()
	if err != nil {
		return nil, err
	}

	return req, nil
}

func NewRequest_ShowBinary(uid string) (req *Request, err error) {
	if !IsUidValid(uid) {
		return nil, fmt.Errorf(ErrUid)
	}

	req = &Request{
		SRS:          0, // Will be automatically calculated.
		RequestSizeA: 0, // Will be automatically calculated.
		RequestSizeB: 0, // Will be automatically calculated.
		Method:       MethodShowBinary,
		UID:          uid,
	}

	err = req.calculateSRS(uid)
	if err != nil {
		return nil, err
	}

	err = req.calculateRequestSize()
	if err != nil {
		return nil, err
	}

	return req, nil
}

func NewRequest_RemoveText(uid string) (req *Request, err error) {
	if !IsUidValid(uid) {
		return nil, fmt.Errorf(ErrUid)
	}

	req = &Request{
		SRS:          0, // Will be automatically calculated.
		RequestSizeA: 0, // Will be automatically calculated.
		RequestSizeB: 0, // Will be automatically calculated.
		Method:       MethodForgetTextRecord,
		UID:          uid,
	}

	err = req.calculateSRS(uid)
	if err != nil {
		return nil, err
	}

	err = req.calculateRequestSize()
	if err != nil {
		return nil, err
	}

	return req, nil
}

func NewRequest_RemoveBinary(uid string) (req *Request, err error) {
	if !IsUidValid(uid) {
		return nil, fmt.Errorf(ErrUid)
	}

	req = &Request{
		SRS:          0, // Will be automatically calculated.
		RequestSizeA: 0, // Will be automatically calculated.
		RequestSizeB: 0, // Will be automatically calculated.
		Method:       MethodForgetBinaryRecord,
		UID:          uid,
	}

	err = req.calculateSRS(uid)
	if err != nil {
		return nil, err
	}

	err = req.calculateRequestSize()
	if err != nil {
		return nil, err
	}

	return req, nil
}

func NewRequest_ClearTextCache() (req *Request, err error) {
	return &Request{
		SRS:          SRS_A,
		RequestSizeA: MethodNameLengthLimit,
		Method:       MethodResetTextCache,
	}, nil
}

func NewRequest_ClearBinaryCache() (req *Request, err error) {
	return &Request{
		SRS:          SRS_A,
		RequestSizeA: MethodNameLengthLimit,
		Method:       MethodResetBinaryCache,
	}, nil
}

func (r *Request) calculateSRS(uid string) (err error) {
	uidLen := len(uid)

	if uidLen <= 0 {
		return errors.New(ErrUid)
	}

	if uidLen <= RequestMessageLengthA-MethodNameLengthLimit {
		r.SRS = SRS_A
		return nil
	}

	if uidLen <= RequestMessageLengthB-MethodNameLengthLimit {
		r.SRS = SRS_B
		return nil
	}

	return errors.New(ErrUidIsTooLong)
}

func (r *Request) calculateRequestSize() (err error) {
	if (r.Method == MethodShowText) || (r.Method == MethodShowBinary) ||
		(r.Method == MethodForgetTextRecord) || (r.Method == MethodForgetBinaryRecord) {
		if r.SRS == SRS_A {
			r.RequestSizeA = MethodNameLengthLimit + uint8(len(r.UID))
		} else if r.SRS == SRS_B {
			r.RequestSizeB = MethodNameLengthLimit + uint16(len(r.UID))
		} else {
			return fmt.Errorf(ErrSrsIsNotSupported, r.SRS)
		}
		return nil
	}

	return fmt.Errorf(ErrUnsupportedMethodValue, r.Method)
}
