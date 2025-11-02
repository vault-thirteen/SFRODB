package connection

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"net"

	ce "github.com/vault-thirteen/SFRODB/pkg/SFRODB/classes/CommonError"
	"github.com/vault-thirteen/SFRODB/pkg/SFRODB/classes/Endianness"
	"github.com/vault-thirteen/SFRODB/pkg/SFRODB/classes/Method"
	"github.com/vault-thirteen/SFRODB/pkg/SFRODB/classes/Request"
	"github.com/vault-thirteen/SFRODB/pkg/SFRODB/classes/Response"
	"github.com/vault-thirteen/SFRODB/pkg/SFRODB/classes/Status"
	uid "github.com/vault-thirteen/SFRODB/pkg/SFRODB/classes/UID"
	"github.com/vault-thirteen/SFRODB/pkg/SFRODB/protocol"
	"github.com/vault-thirteen/SFRODB/pkg/SFRODB/std/tcp"
)

type Connection struct {
	netConn                    *net.TCPConn
	responseMessageLengthLimit uint
	clientId                   string
}

func New(
	netConn *net.TCPConn,
	responseMessageLengthLimit uint,
	clientId string,
) (con *Connection) {
	return &Connection{
		netConn:                    netConn,
		responseMessageLengthLimit: responseMessageLengthLimit,
		clientId:                   clientId,
	}
}

func (con *Connection) ClientId() (clientId string) {
	return con.clientId
}

// Break is a method used by a Client to finalise its connection.
func (con *Connection) Break() (cerr *ce.CommonError) {
	err := con.netConn.Close()
	if err != nil {
		return ce.NewServerError(err.Error(), 0, 0, con.clientId)
	}

	return nil
}

// SendRequestMessage is a method used by a Client to send a request to the
// server.
func (con *Connection) SendRequestMessage(req *request.Request) (cerr *ce.CommonError) {
	var buf bytes.Buffer
	var err error
	var ba []byte

	// 1. Size.
	{
		if req.Size == 0 {
			rs := protocol.MethodNameLen + req.UID.Length()
			if rs > math.MaxUint16 {
				return ce.NewClientError(fmt.Sprintf(request.ErrSizeIsTooLong, rs), req.Method, 0, con.clientId)
			}
			req.Size = uint16(rs)
		}

		ba = make([]byte, protocol.RequestSizeLen)
		switch protocol.Endianness {
		case endianness.Endianness_BigEndian:
			binary.BigEndian.PutUint16(ba, req.Size)

		case endianness.Endianness_LittleEndian:
			binary.LittleEndian.PutUint16(ba, req.Size)

		default:
			return ce.NewClientError(endianness.ErrEndiannessIsUnknown, req.Method, 0, con.clientId)
		}

		_, err = buf.Write(ba)
		if err != nil {
			return ce.NewClientError(err.Error(), req.Method, 0, con.clientId)
		}
	}

	// 2. Method.
	{
		ba, err = req.Method.Bytes()
		if err != nil {
			return ce.NewClientError(err.Error(), req.Method, 0, con.clientId)
		}

		_, err = buf.Write(ba)
		if err != nil {
			return ce.NewClientError(err.Error(), req.Method, 0, con.clientId)
		}
	}

	// 3. UID.
	{
		if req.UID != nil {
			_, err = buf.Write(req.UID.Bytes())
			if err != nil {
				return ce.NewClientError(err.Error(), req.Method, 0, con.clientId)
			}
		}
	}

	// Send data.
	_, err = con.netConn.Write(buf.Bytes())
	if err != nil {
		return ce.NewClientError(err.Error(), req.Method, 0, con.clientId)
	}

	return nil
}

// GetResponseMessage is a method used by a Client to read a response from the
// server.
func (con *Connection) GetResponseMessage() (resp *response.Response, cerr *ce.CommonError) {
	resp = &response.Response{}
	var err error
	var ba []byte

	// 1. Size.
	{
		ba, err = tcp.ReadExactSize(con.netConn, protocol.ResponseSizeLen)
		if err != nil {
			return nil, ce.NewClientError(err.Error(), 0, 0, con.clientId)
		}

		switch protocol.Endianness {
		case endianness.Endianness_BigEndian:
			resp.Size = binary.BigEndian.Uint32(ba)

		case endianness.Endianness_LittleEndian:
			resp.Size = binary.LittleEndian.Uint32(ba)

		default:
			return nil, ce.NewClientError(endianness.ErrEndiannessIsUnknown, 0, 0, con.clientId)
		}

		if resp.Size < protocol.StatusNameLen {
			return nil, ce.NewClientError(fmt.Sprintf(response.ErrSizeIsTooShort, resp.Size), 0, 0, con.clientId)
		}
	}

	// 2. Status.
	{
		ba, err = tcp.ReadExactSize(con.netConn, protocol.StatusNameLen)
		if err != nil {
			return nil, ce.NewClientError(err.Error(), 0, 0, con.clientId)
		}

		resp.Status, err = status.NewFromString(string(ba[0:protocol.StatusNameLen]))
		if err != nil {
			return nil, ce.NewClientError(err.Error(), 0, 0, con.clientId)
		}
	}

	// 3. Data.
	{
		dataSize := uint(resp.Size) - uint(protocol.StatusNameLen)
		if dataSize > 0 {
			resp.Data, err = tcp.ReadExactSize(con.netConn, dataSize)
			if err != nil {
				return nil, ce.NewClientError(err.Error(), 0, resp.Status, con.clientId)
			}
		}
	}

	return resp, nil
}

// GetNextRequest is a method used by a Server to receive a request from the
// client.
func (con *Connection) GetNextRequest() (req *request.Request, cerr *ce.CommonError) {
	req = &request.Request{}
	var err error
	var ba []byte

	// 1. Size.
	{
		ba, err = tcp.ReadExactSize(con.netConn, protocol.RequestSizeLen)
		if err != nil {
			return nil, ce.NewServerError(err.Error(), 0, 0, con.clientId)
		}

		switch protocol.Endianness {
		case endianness.Endianness_BigEndian:
			req.Size = binary.BigEndian.Uint16(ba)

		case endianness.Endianness_LittleEndian:
			req.Size = binary.LittleEndian.Uint16(ba)

		default:
			return nil, ce.NewServerError(endianness.ErrEndiannessIsUnknown, 0, 0, con.clientId)
		}

		if req.Size < protocol.MethodNameLen {
			return nil, ce.NewServerError(fmt.Sprintf(request.ErrSizeIsTooShort, req.Size), 0, 0, con.clientId)
		}
	}

	// 2. Method.
	{
		ba, err = tcp.ReadExactSize(con.netConn, protocol.MethodNameLen)
		if err != nil {
			return nil, ce.NewServerError(err.Error(), 0, 0, con.clientId)
		}

		req.Method, err = method.NewFromString(string(ba[0:protocol.MethodNameLen]))
		if err != nil {
			return nil, ce.NewServerError(err.Error(), 0, 0, con.clientId)
		}
	}

	// 3. UID.
	{
		uidSize := uint(req.Size) - uint(protocol.MethodNameLen)
		if uidSize > 0 {
			ba, err = tcp.ReadExactSize(con.netConn, uidSize)
			if err != nil {
				return nil, ce.NewServerError(err.Error(), req.Method, 0, con.clientId)
			}
		}

		req.UID, err = uid.New(string(ba[0:uidSize]))
		if err != nil {
			return nil, ce.NewServerError(err.Error(), req.Method, 0, con.clientId)
		}
	}

	return req, nil
}

// SendResponseMessage is a method used by a Server to send a response to the
// client.
func (con *Connection) SendResponseMessage(resp *response.Response) (cerr *ce.CommonError) {
	var buf bytes.Buffer
	var ba []byte
	var err error

	// 1. Size.
	{
		if resp.Size == 0 {
			rs := protocol.StatusNameLen + len(resp.Data)
			if rs > math.MaxUint32 {
				return ce.NewServerError(fmt.Sprintf(response.ErrSizeIsTooLong, rs), 0, resp.Status, con.clientId)
			}
			resp.Size = uint32(rs)
		}

		ba = make([]byte, protocol.ResponseSizeLen)
		switch protocol.Endianness {
		case endianness.Endianness_BigEndian:
			binary.BigEndian.PutUint32(ba, resp.Size)

		case endianness.Endianness_LittleEndian:
			binary.LittleEndian.PutUint32(ba, resp.Size)

		default:
			return ce.NewServerError(endianness.ErrEndiannessIsUnknown, 0, resp.Status, con.clientId)
		}

		_, err = buf.Write(ba)
		if err != nil {
			return ce.NewServerError(err.Error(), 0, resp.Status, con.clientId)
		}
	}

	// 2. Status.
	{
		ba, err = resp.Status.Bytes()
		if err != nil {
			return ce.NewServerError(err.Error(), 0, resp.Status, con.clientId)
		}

		_, err = buf.Write(ba)
		if err != nil {
			return ce.NewServerError(err.Error(), 0, resp.Status, con.clientId)
		}
	}

	// 3. Data.
	{
		_, err = buf.Write(resp.Data)
		if err != nil {
			return ce.NewServerError(err.Error(), 0, resp.Status, con.clientId)
		}
	}

	// Send data.
	_, err = con.netConn.Write(buf.Bytes())
	if err != nil {
		return ce.NewServerError(err.Error(), 0, resp.Status, con.clientId)
	}

	return nil
}
