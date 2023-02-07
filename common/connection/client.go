package connection

import (
	ce "github.com/vault-thirteen/SFRODB/common/error"
	"github.com/vault-thirteen/SFRODB/common/request"
	"github.com/vault-thirteen/SFRODB/common/response"
)

// Break is a method used by a Client to finalize its connection.
func (c *Connection) Break() (cerr *ce.CommonError) {
	err := c.close()
	if err != nil {
		return ce.NewServerError(err.Error(), 0)
	}

	return nil
}

// SendRequestMessage is a method used by a Client to send a request to the
// server.
func (c *Connection) SendRequestMessage(rm *request.Request) (cerr *ce.CommonError) {
	var err error
	err = c.sendSRS(rm.SRS)
	if err != nil {
		return ce.NewServerError(err.Error(), 0)
	}

	err = c.sendRequestSize(rm)
	if err != nil {
		return ce.NewServerError(err.Error(), 0)
	}

	err = c.sendRequestMethodAndUid(rm)
	if err != nil {
		return ce.NewServerError(err.Error(), 0)
	}

	return nil
}

// GetResponseMessage is a method used by a Client to read a response from the
// server.
func (c *Connection) GetResponseMessage(useBinary bool) (resp *response.Response, cerr *ce.CommonError) {
	resp = &response.Response{}

	var err error
	resp.SRS, err = c.getSRS()
	if err != nil {
		return nil, ce.NewServerError(ce.ErrSrsReading+err.Error(), 0)
	}

	err = c.getResponseSize(resp)
	if err != nil {
		return nil, ce.NewServerError(ce.ErrRsReading+err.Error(), 0)
	}

	err = c.getResponseMethodAndData(resp, useBinary)
	if err != nil {
		return nil, ce.NewServerError(ce.ErrReadingMethodAndData+err.Error(), 0)
	}

	return resp, nil
}
