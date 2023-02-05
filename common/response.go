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

func (r *Response) IsClosingConnection() bool {
	return r.Method == MethodClosingConnection
}

func (r *Response) IsShowingText() bool {
	return r.Method == MethodShowingText
}

func (r *Response) IsShowingBinary() bool {
	return r.Method == MethodShowingBinary
}

func NewResponse_ClosingConnection() (resp *Response, err error) {
	return &Response{
		SRS:           SRS_A,
		ResponseSizeA: MethodNameLengthLimit,
		Method:        MethodClosingConnection,
	}, nil
}

func NewResponse_ClientErrorWarning() (resp *Response, err error) {
	return &Response{
		SRS:           SRS_A,
		ResponseSizeA: MethodNameLengthLimit,
		Method:        MethodClientError,
	}, nil
}

func NewResponse_ShowingText(text string) (resp *Response, err error) {
	resp = &Response{
		SRS:           0, // Will be automatically calculated.
		ResponseSizeC: 0, // Will be automatically calculated.
		Method:        MethodShowingText,
		Text:          text,
	}

	if len(text) > ResponseMessageLengthC-MethodNameLengthLimit {
		err = fmt.Errorf(ErrTextIsTooLong, ResponseMessageLengthC-MethodNameLengthLimit, len(text))
		return nil, err
	} else if len(text) > ResponseMessageLengthB-MethodNameLengthLimit {
		resp.SRS = SRS_C
	} else if len(text) > ResponseMessageLengthA-MethodNameLengthLimit {
		resp.SRS = SRS_B
	} else {
		resp.SRS = SRS_A
	}

	err = resp.calculateResponseSize()
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func NewResponse_ShowingBinary(data []byte) (resp *Response, err error) {
	resp = &Response{
		SRS:           0, // Will be automatically calculated.
		ResponseSizeC: 0, // Will be automatically calculated.
		Method:        MethodShowingBinary,
		Data:          data,
	}

	if len(data) > ResponseMessageLengthC-MethodNameLengthLimit {
		err = fmt.Errorf(ErrTextIsTooLong, ResponseMessageLengthC-MethodNameLengthLimit, len(data))
		return nil, err
	} else if len(data) > ResponseMessageLengthB-MethodNameLengthLimit {
		resp.SRS = SRS_C
	} else if len(data) > ResponseMessageLengthA-MethodNameLengthLimit {
		resp.SRS = SRS_B
	} else {
		resp.SRS = SRS_A
	}

	err = resp.calculateResponseSize()
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (r *Response) calculateResponseSize() (err error) {
	switch r.Method {
	case MethodClosingConnection:
		r.ResponseSizeA = MethodNameLengthLimit

	case MethodShowingText:
		switch r.SRS {
		case SRS_A:
			r.ResponseSizeA = MethodNameLengthLimit + uint8(len(r.Text))
		case SRS_B:
			r.ResponseSizeB = MethodNameLengthLimit + uint16(len(r.Text))
		case SRS_C:
			r.ResponseSizeC = MethodNameLengthLimit + uint32(len(r.Text))
		default:
			return fmt.Errorf(ErrSrsIsNotSupported, r.SRS)
		}

	case MethodShowingBinary:
		switch r.SRS {
		case SRS_A:
			r.ResponseSizeA = MethodNameLengthLimit + uint8(len(r.Data))
		case SRS_B:
			r.ResponseSizeB = MethodNameLengthLimit + uint16(len(r.Data))
		case SRS_C:
			r.ResponseSizeC = MethodNameLengthLimit + uint32(len(r.Data))
		default:
			return fmt.Errorf(ErrSrsIsNotSupported, r.SRS)
		}

	default:
		return fmt.Errorf(ErrUnsupportedMethodValue, r.Method)
	}

	return nil
}

func NewResponse_OK() (resp *Response, err error) {
	return &Response{
		SRS:           SRS_A,
		ResponseSizeA: MethodNameLengthLimit,
		Method:        MethodOK,
	}, nil
}
