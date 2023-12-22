package connection

import (
	"bytes"

	"github.com/vault-thirteen/SFRODB/pkg/SFRODB/common/error"
	"github.com/vault-thirteen/SFRODB/pkg/SFRODB/common/request"
	"github.com/vault-thirteen/SFRODB/pkg/SFRODB/common/response"
)

// GetNextRequest is a method used by a Server to receive a request from the
// client.
func (con *Connection) GetNextRequest() (r *request.Request, cerr *ce.CommonError) {
	r = &request.Request{}

	var err error
	r.SRS, err = con.getSRS()
	if err != nil {
		return nil, ce.NewServerError(ce.ErrSrsReading+err.Error(), 0, con.clientId)
	}

	err = con.getRequestSize(r)
	if err != nil {
		return nil, ce.NewServerError(ce.ErrRsReading+err.Error(), 0, con.clientId)
	}

	err = con.getRequestMethodAndUID(r)
	if err != nil {
		return nil, ce.NewServerError(ce.ErrReadingMethodAndData+err.Error(), 0, con.clientId)
	}

	return r, nil
}

// SendResponseMessage is a method used by a Server to send a response to the
// client.
func (con *Connection) SendResponseMessage(rm *response.Response) (cerr *ce.CommonError) {
	var buf bytes.Buffer
	var ba []byte
	var err error
	ba = con.writeSRS(rm.SRS)
	_, err = buf.Write(ba)
	if err != nil {
		return ce.NewServerError(err.Error(), 0, con.clientId)
	}

	ba = con.writeResponseSize(rm)
	_, err = buf.Write(ba)
	if err != nil {
		return ce.NewServerError(err.Error(), 0, con.clientId)
	}

	ba, err = con.writeResponseMethodAndData(rm)
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
