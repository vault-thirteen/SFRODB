package common

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/vault-thirteen/SFRODB/common/reader"
)

type Connection struct {
	netConn                    net.Conn
	methodNameBuffers          *map[Method][]byte
	methodValues               *map[string]Method
	responseMessageLengthLimit uint
}

func NewConnection(
	netConn net.Conn,
	methodNameBuffers *map[Method][]byte,
	methodValues *map[string]Method,

	// This limit is used on client only.
	responseMessageLengthLimit uint,
) (c *Connection, err error) {
	return &Connection{
		netConn:                    netConn,
		methodNameBuffers:          methodNameBuffers,
		methodValues:               methodValues,
		responseMessageLengthLimit: responseMessageLengthLimit,
	}, nil
}

// Finalize is a method used by a Server to finalize the client's connection.
// This method is used either when the client requested to stop the
// communication or when an internal error happened on the server.
func (c *Connection) Finalize() (err error) {
	var rm *Response
	rm, err = NewResponse_ClosingConnection()
	if err != nil {
		return err
	}

	err = c.SendResponseMessage(rm)
	if err != nil {
		return err
	}

	return c.close()
}

// Warn is a method used by a Server to warn the client about its (client's)
// error.
func (c *Connection) Warn() (err error) {
	var rm *Response
	rm, err = NewResponse_ClientErrorWarning()
	if err != nil {
		return err
	}

	err = c.SendResponseMessage(rm)
	if err != nil {
		return err
	}

	return nil
}

// Break is a method used by a Client to finalize its connection.
func (c *Connection) Break() (err error) {
	return c.close()
}

func (c *Connection) close() (err error) {
	return c.netConn.Close()
}

// GetNextRequest is a method used by a Server to receive a request from the
// client.
func (c *Connection) GetNextRequest() (r *Request, err error) {
	r = &Request{}

	r.SRS, err = c.getSRS()
	if err != nil {
		return nil, errors.New(ErrSrsReading + err.Error())
	}

	err = c.getRequestSize(r)
	if err != nil {
		return nil, errors.New(ErrRsReading + err.Error())
	}

	err = c.getRequestMethodAndUID(r)
	if err != nil {
		return nil, errors.New(ErrReadingMethodAndData + err.Error())
	}

	return r, nil
}

func (c *Connection) getSRS() (srs byte, err error) {
	var data []byte
	data, err = reader.ReadExactSize(c.netConn, 1)
	if err != nil {
		return 0, err
	}

	srs = data[0]

	switch srs {
	case SRS_A:
		return srs, nil
	case SRS_B:
		return srs, nil
	case SRS_C:
		return srs, nil
	}

	return 0, fmt.Errorf(ErrSrsIsNotSupported, srs)
}

func (c *Connection) getRequestSize(r *Request) (err error) {
	switch r.SRS {
	case SRS_A:
		return c.getRequestSizeA(r)
	case SRS_B:
		return c.getRequestSizeB(r)
	}

	return fmt.Errorf(ErrSrsIsNotSupported, r.SRS)
}

func (c *Connection) getRequestSizeA(r *Request) (err error) {
	var data []byte
	data, err = reader.ReadExactSize(c.netConn, RS_LengthA)
	if err != nil {
		return err
	}

	r.RequestSizeA = data[0]

	return nil
}

func (c *Connection) getRequestSizeB(r *Request) (err error) {
	var data []byte
	data, err = reader.ReadExactSize(c.netConn, RS_LengthB)
	if err != nil {
		return err
	}

	r.RequestSizeB = binary.BigEndian.Uint16(data)

	return nil
}

func (c *Connection) getRequestMethodAndUID(r *Request) (err error) {
	var reqMsgLen uint
	switch r.SRS {
	case SRS_A:
		reqMsgLen = uint(r.RequestSizeA)
	case SRS_B:
		reqMsgLen = uint(r.RequestSizeB)
	default:
		return fmt.Errorf(ErrSrsIsNotSupported, r.SRS)
	}

	var data []byte
	data, err = reader.ReadExactSize(c.netConn, reqMsgLen)
	if err != nil {
		return err
	}

	r.Method, err = c.NewMethodFromBytes(data[0:3])
	if err != nil {
		return err
	}

	r.UID = strings.TrimSpace(string(data[3:reqMsgLen]))

	return nil
}

// SendResponseMessage is a method used by a Server to send a response to the
// client.
func (c *Connection) SendResponseMessage(rm *Response) (err error) {
	err = c.sendSRS(rm.SRS)
	if err != nil {
		return err
	}

	err = c.sendResponseSize(rm)
	if err != nil {
		return err
	}

	err = c.sendResponseMethodAndData(rm)
	if err != nil {
		return err
	}

	return nil
}

func (c *Connection) sendSRS(srs byte) (err error) {
	buf := make([]byte, 1)
	buf[0] = srs

	_, err = c.netConn.Write(buf)
	if err != nil {
		return err
	}

	return nil
}

func (c *Connection) sendResponseSize(rm *Response) (err error) {
	var buf []byte
	switch rm.SRS {
	case SRS_A:
		buf = make([]byte, RS_LengthA)
		buf[0] = rm.ResponseSizeA

	case SRS_B:
		buf = make([]byte, RS_LengthB)
		binary.BigEndian.PutUint16(buf, rm.ResponseSizeB)

	case SRS_C:
		buf = make([]byte, RS_LengthC)
		binary.BigEndian.PutUint32(buf, rm.ResponseSizeC)
	}

	_, err = c.netConn.Write(buf)
	if err != nil {
		return err
	}

	return nil
}

func (c *Connection) sendResponseMethodAndData(rm *Response) (err error) {
	_, err = c.netConn.Write((*c.methodNameBuffers)[rm.Method])
	if err != nil {
		return err
	}

	switch rm.Method {
	case MethodShowingText:
		_, err = c.netConn.Write([]byte(rm.Text))
	case MethodShowingBinary:
		_, err = c.netConn.Write(rm.Data)
	}
	if err != nil {
		return err
	}

	return nil
}

// SendRequestMessage is a method used by a Client to send a request to the
// server.
func (c *Connection) SendRequestMessage(rm *Request) (err error) {
	err = c.sendSRS(rm.SRS)
	if err != nil {
		return err
	}

	err = c.sendRequestSize(rm)
	if err != nil {
		return err
	}

	err = c.sendRequestMethodAndUid(rm)
	if err != nil {
		return err
	}

	return nil
}

func (c *Connection) sendRequestSize(rm *Request) (err error) {
	var buf []byte
	switch rm.SRS {
	case SRS_A:
		buf = make([]byte, RS_LengthA)
		buf[0] = rm.RequestSizeA

	case SRS_B:
		buf = make([]byte, RS_LengthB)
		binary.BigEndian.PutUint16(buf, rm.RequestSizeB)

	default:
		return fmt.Errorf(ErrSrsIsNotSupported, rm.SRS)
	}

	_, err = c.netConn.Write(buf)
	if err != nil {
		return err
	}

	return nil
}

func (c *Connection) sendRequestMethodAndUid(rm *Request) (err error) {
	switch rm.Method {
	case MethodShowText,
		MethodShowBinary,
		MethodForgetTextRecord,
		MethodForgetBinaryRecord,
		MethodResetTextCache,
		MethodResetBinaryCache,
		MethodCloseConnection:
		break
	default:
		return fmt.Errorf(ErrUnsupportedMethodValue, rm.Method)
	}

	_, err = c.netConn.Write((*c.methodNameBuffers)[rm.Method])
	if err != nil {
		return err
	}

	_, err = c.netConn.Write([]byte(rm.UID))
	if err != nil {
		return err
	}

	return nil
}

// GetResponseMessage is a method used by a Client to read a response from the
// server.
func (c *Connection) GetResponseMessage() (resp *Response, err error) {
	resp = &Response{}

	resp.SRS, err = c.getSRS()
	if err != nil {
		return nil, errors.New(ErrSrsReading + err.Error())
	}

	err = c.getResponseSize(resp)
	if err != nil {
		return nil, errors.New(ErrRsReading + err.Error())
	}

	err = c.getResponseMethodAndData(resp)
	if err != nil {
		return nil, errors.New(ErrReadingMethodAndData + err.Error())
	}

	return resp, nil
}

func (c *Connection) getResponseSize(resp *Response) (err error) {
	switch resp.SRS {
	case SRS_A:
		return c.getResponseSizeA(resp)
	case SRS_B:
		return c.getResponseSizeB(resp)
	case SRS_C:
		return c.getResponseSizeC(resp)
	}

	return fmt.Errorf(ErrSrsIsNotSupported, resp.SRS)
}

func (c *Connection) getResponseSizeA(resp *Response) (err error) {
	var data []byte
	data, err = reader.ReadExactSize(c.netConn, RS_LengthA)
	if err != nil {
		return err
	}

	resp.ResponseSizeA = data[0]

	return nil
}

func (c *Connection) getResponseSizeB(resp *Response) (err error) {
	var data []byte
	data, err = reader.ReadExactSize(c.netConn, RS_LengthB)
	if err != nil {
		return err
	}

	resp.ResponseSizeB = binary.BigEndian.Uint16(data)

	return nil
}

func (c *Connection) getResponseSizeC(resp *Response) (err error) {
	var data []byte
	data, err = reader.ReadExactSize(c.netConn, RS_LengthC)
	if err != nil {
		return err
	}

	resp.ResponseSizeC = binary.BigEndian.Uint32(data)

	return nil
}

func (c *Connection) getResponseMethodAndData(resp *Response) (err error) {
	var respMsgLen uint
	switch resp.SRS {
	case SRS_A:
		respMsgLen = uint(resp.ResponseSizeA)
	case SRS_B:
		respMsgLen = uint(resp.ResponseSizeB)
	case SRS_C:
		respMsgLen = uint(resp.ResponseSizeC)
	default:
		return fmt.Errorf(ErrSrsIsNotSupported, resp.SRS)
	}

	if respMsgLen > c.responseMessageLengthLimit {
		return fmt.Errorf(ErrMessageIsTooLong, c.responseMessageLengthLimit, respMsgLen)
	}

	var data []byte
	data, err = reader.ReadExactSize(c.netConn, respMsgLen)
	if err != nil {
		return err
	}

	resp.Method, err = c.NewMethodFromBytes(data[0:3])
	if err != nil {
		return err
	}

	switch resp.Method {
	case MethodShowingText:
		resp.Text = strings.TrimSpace(string(data[3:respMsgLen]))
	case MethodShowingBinary:
		resp.Data = data[3:respMsgLen]
	}

	return nil
}
