package connection

import (
	"bytes"

	"github.com/vault-thirteen/SFRODB/pkg/common/error"
	"github.com/vault-thirteen/SFRODB/pkg/common/request"
	"github.com/vault-thirteen/SFRODB/pkg/common/response"
)

// Break is a method used by a Client to finalize its connection.
func (con *Connection) Break() (cerr *ce.CommonError) {
	err := con.close()
	if err != nil {
		return ce.NewServerError(err.Error(), 0, con.clientId)
	}

	return nil
}

// SendRequestMessage is a method used by a Client to send a request to the
// server.
func (con *Connection) SendRequestMessage(rm *request.Request) (cerr *ce.CommonError) {
	var buf bytes.Buffer
	var ba []byte
	var err error
	ba = con.writeSRS(rm.SRS)
	_, err = buf.Write(ba)
	if err != nil {
		return ce.NewServerError(err.Error(), 0, con.clientId)
	}

	ba = con.writeRequestSize(rm)
	_, err = buf.Write(ba)
	if err != nil {
		return ce.NewServerError(err.Error(), 0, con.clientId)
	}

	ba, err = con.writeRequestMethodAndUid(rm)
	if err != nil {
		return ce.NewServerError(err.Error(), 0, con.clientId)
	}
	_, err = buf.Write(ba)
	if err != nil {
		return ce.NewServerError(err.Error(), 0, con.clientId)
	}

	err = con.sendData(buf.Bytes())
	if err != nil {
		return ce.NewServerError(err.Error(), 0, con.clientId)
	}

	return nil
}

// GetResponseMessage is a method used by a Client to read a response from the
// server.
func (con *Connection) GetResponseMessage() (resp *response.Response, cerr *ce.CommonError) {
	resp = &response.Response{}

	var err error
	resp.SRS, err = con.getSRS()
	if err != nil {
		return nil, ce.NewServerError(ce.ErrSrsReading+err.Error(), 0, con.clientId)
	}

	err = con.getResponseSize(resp)
	if err != nil {
		return nil, ce.NewServerError(ce.ErrRsReading+err.Error(), 0, con.clientId)
	}

	err = con.getResponseMethodAndData(resp)
	if err != nil {
		return nil, ce.NewServerError(ce.ErrReadingMethodAndData+err.Error(), 0, con.clientId)
	}

	return resp, nil
}
