package connection

import (
	ce "github.com/vault-thirteen/SFRODB/common/error"
	"github.com/vault-thirteen/SFRODB/common/request"
	"github.com/vault-thirteen/SFRODB/common/response"
)

// GetNextRequest is a method used by a Server to receive a request from the
// client.
func (c *Connection) GetNextRequest() (r *request.Request, cerr *ce.CommonError) {
	r = &request.Request{}

	var err error
	r.SRS, err = c.getSRS()
	if err != nil {
		return nil, ce.NewServerError(ce.ErrSrsReading+err.Error(), 0)
	}

	err = c.getRequestSize(r)
	if err != nil {
		return nil, ce.NewServerError(ce.ErrRsReading+err.Error(), 0)
	}

	err = c.getRequestMethodAndUID(r)
	if err != nil {
		return nil, ce.NewServerError(ce.ErrReadingMethodAndData+err.Error(), 0)
	}

	return r, nil
}

// SendResponseMessage is a method used by a Server to send a response to the
// client.
func (c *Connection) SendResponseMessage(rm *response.Response, useBinary bool) (cerr *ce.CommonError) {
	var err error
	err = c.sendSRS(rm.SRS)
	if err != nil {
		return ce.NewServerError(err.Error(), 0)
	}

	err = c.sendResponseSize(rm)
	if err != nil {
		return ce.NewServerError(err.Error(), 0)
	}

	err = c.sendResponseMethodAndData(rm, useBinary)
	if err != nil {
		return ce.NewServerError(err.Error(), 0)
	}

	return nil
}
