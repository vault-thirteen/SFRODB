package connection

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/vault-thirteen/SFRODB/pkg/SFRODB/common/error"
	"github.com/vault-thirteen/SFRODB/pkg/SFRODB/common/method"
	mn "github.com/vault-thirteen/SFRODB/pkg/SFRODB/common/method/name"
	proto "github.com/vault-thirteen/SFRODB/pkg/SFRODB/common/protocol"
	"github.com/vault-thirteen/SFRODB/pkg/SFRODB/common/reader"
	"github.com/vault-thirteen/SFRODB/pkg/SFRODB/common/request"
	"github.com/vault-thirteen/SFRODB/pkg/SFRODB/common/response"
)

type Connection struct {
	netConn                    *net.TCPConn
	methodNameBuffers          *map[method.Method][]byte
	methodValues               *map[string]method.Method
	responseMessageLengthLimit uint
	clientId                   string
}

func NewConnection(
	netConn *net.TCPConn,
	methodNameBuffers *map[method.Method][]byte,
	methodValues *map[string]method.Method,

	// This limit is used on client only.
	responseMessageLengthLimit uint,
	clientId string,
) (con *Connection) {
	return &Connection{
		netConn:                    netConn,
		methodNameBuffers:          methodNameBuffers,
		methodValues:               methodValues,
		responseMessageLengthLimit: responseMessageLengthLimit,
		clientId:                   clientId,
	}
}

func (con *Connection) close() (err error) {
	return con.netConn.Close()
}

func (con *Connection) getSRS() (srs byte, err error) {
	var data []byte
	data, err = reader.ReadExactSize(con.netConn, 1)
	if err != nil {
		return 0, err
	}

	srs = data[0]

	switch srs {
	case proto.SRS_A:
		return srs, nil
	case proto.SRS_B:
		return srs, nil
	case proto.SRS_C:
		return srs, nil
	}

	return 0, fmt.Errorf(ce.ErrSrsIsNotSupported, srs)
}

func (con *Connection) getRequestSize(r *request.Request) (err error) {
	switch r.SRS {
	case proto.SRS_A:
		return con.getRequestSizeA(r)
	case proto.SRS_B:
		return con.getRequestSizeB(r)
	}

	return fmt.Errorf(ce.ErrSrsIsNotSupported, r.SRS)
}

func (con *Connection) getRequestSizeA(r *request.Request) (err error) {
	var data []byte
	data, err = reader.ReadExactSize(con.netConn, proto.RS_LengthA)
	if err != nil {
		return err
	}

	r.RequestSizeA = data[0]

	return nil
}

func (con *Connection) getRequestSizeB(r *request.Request) (err error) {
	var data []byte
	data, err = reader.ReadExactSize(con.netConn, proto.RS_LengthB)
	if err != nil {
		return err
	}

	r.RequestSizeB = binary.BigEndian.Uint16(data)

	return nil
}

func (con *Connection) getRequestMethodAndUID(r *request.Request) (err error) {
	var reqMsgLen uint
	switch r.SRS {
	case proto.SRS_A:
		reqMsgLen = uint(r.RequestSizeA)
	case proto.SRS_B:
		reqMsgLen = uint(r.RequestSizeB)
	default:
		return fmt.Errorf(ce.ErrSrsIsNotSupported, r.SRS)
	}

	var data []byte
	data, err = reader.ReadExactSize(con.netConn, reqMsgLen)
	if err != nil {
		return err
	}

	r.Method, err = con.NewMethodFromBytes(data[0:3])
	if err != nil {
		return err
	}

	r.UID = strings.TrimSpace(string(data[3:reqMsgLen]))

	return nil
}

func (con *Connection) writeSRS(srs byte) (ba []byte) {
	return []byte{srs}
}

func (con *Connection) writeResponseSize(rm *response.Response) (ba []byte) {
	switch rm.SRS {
	case proto.SRS_A:
		ba = make([]byte, proto.RS_LengthA)
		ba[0] = rm.ResponseSizeA
		return ba

	case proto.SRS_B:
		ba = make([]byte, proto.RS_LengthB)
		binary.BigEndian.PutUint16(ba, rm.ResponseSizeB)
		return ba

	case proto.SRS_C:
		ba = make([]byte, proto.RS_LengthC)
		binary.BigEndian.PutUint32(ba, rm.ResponseSizeC)
		return ba

	default:
		return nil
	}
}

func (con *Connection) writeResponseMethodAndData(rm *response.Response) (ba []byte, err error) {
	var buf bytes.Buffer
	_, err = buf.Write((*con.methodNameBuffers)[rm.Method])
	if err != nil {
		return nil, err
	}

	_, err = buf.Write(rm.Data)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (con *Connection) writeRequestSize(rm *request.Request) (ba []byte) {
	switch rm.SRS {
	case proto.SRS_A:
		ba = make([]byte, proto.RS_LengthA)
		ba[0] = rm.RequestSizeA
		return ba

	case proto.SRS_B:
		ba = make([]byte, proto.RS_LengthB)
		binary.BigEndian.PutUint16(ba, rm.RequestSizeB)
		return ba

	default:
		return nil
	}
}

func (con *Connection) writeRequestMethodAndUid(rm *request.Request) (ba []byte, err error) {
	var buf bytes.Buffer
	_, err = buf.Write((*con.methodNameBuffers)[rm.Method])
	if err != nil {
		return nil, err
	}

	_, err = buf.Write([]byte(rm.UID))
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (con *Connection) getResponseSize(resp *response.Response) (err error) {
	switch resp.SRS {
	case proto.SRS_A:
		return con.getResponseSizeA(resp)
	case proto.SRS_B:
		return con.getResponseSizeB(resp)
	case proto.SRS_C:
		return con.getResponseSizeC(resp)
	}

	return fmt.Errorf(ce.ErrSrsIsNotSupported, resp.SRS)
}

func (con *Connection) getResponseSizeA(resp *response.Response) (err error) {
	var data []byte
	data, err = reader.ReadExactSize(con.netConn, proto.RS_LengthA)
	if err != nil {
		return err
	}

	resp.ResponseSizeA = data[0]

	return nil
}

func (con *Connection) getResponseSizeB(resp *response.Response) (err error) {
	var data []byte
	data, err = reader.ReadExactSize(con.netConn, proto.RS_LengthB)
	if err != nil {
		return err
	}

	resp.ResponseSizeB = binary.BigEndian.Uint16(data)

	return nil
}

func (con *Connection) getResponseSizeC(resp *response.Response) (err error) {
	var data []byte
	data, err = reader.ReadExactSize(con.netConn, proto.RS_LengthC)
	if err != nil {
		return err
	}

	resp.ResponseSizeC = binary.BigEndian.Uint32(data)

	return nil
}

func (con *Connection) getResponseMethodAndData(resp *response.Response) (err error) {
	var respMsgLen uint
	switch resp.SRS {
	case proto.SRS_A:
		respMsgLen = uint(resp.ResponseSizeA)
	case proto.SRS_B:
		respMsgLen = uint(resp.ResponseSizeB)
	case proto.SRS_C:
		respMsgLen = uint(resp.ResponseSizeC)
	default:
		return fmt.Errorf(ce.ErrSrsIsNotSupported, resp.SRS)
	}

	if respMsgLen > con.responseMessageLengthLimit {
		return fmt.Errorf(ce.ErrMessageIsTooLong, con.responseMessageLengthLimit, respMsgLen)
	}

	var data []byte
	data, err = reader.ReadExactSize(con.netConn, respMsgLen)
	if err != nil {
		return err
	}

	resp.Method, err = con.NewMethodFromBytes(data[0:3])
	if err != nil {
		return err
	}

	resp.Data = data[3:respMsgLen]

	return nil
}

func (con *Connection) ClientId() (clientId string) {
	return con.clientId
}

func (con *Connection) NewMethodFromBytes(b []byte) (m method.Method, err error) {
	if len(b) == 3 {
		return con.NewMethodFromString(string(b))
	}

	return con.NewMethodFromString(string(b[0:3]))
}

func (con *Connection) NewMethodFromString(s string) (m method.Method, err error) {
	methodStr := strings.TrimSuffix(s, mn.Spacer)

	var ok bool
	m, ok = (*con.methodValues)[methodStr]
	if !ok {
		return 0, fmt.Errorf(ce.ErrUnknownMethodName, methodStr)
	}

	return m, nil
}

func EnableKeepAlives(conn *net.TCPConn) (err error) {
	err = conn.SetKeepAlivePeriod(time.Second * proto.TcpKeepAlivePeriodSec)
	if err != nil {
		return err
	}

	err = conn.SetKeepAlive(proto.TcpKeepAliveIsEnabled)
	if err != nil {
		return err
	}

	return nil
}

func (con *Connection) sendData(data []byte) (err error) {
	_, err = con.netConn.Write(data)
	if err != nil {
		return err
	}

	return nil
}
