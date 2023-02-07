package server

import (
	"github.com/vault-thirteen/SFRODB/common/connection"
	ce "github.com/vault-thirteen/SFRODB/common/error"
	"github.com/vault-thirteen/SFRODB/common/response"
)

// clientError tells the client about its (client's) error.
// Returns a detailed error.
func (srv *Server) clientError(c *connection.Connection) (cerr *ce.CommonError) {
	rm, err := response.NewResponse_ClientError()
	if err != nil {
		return ce.NewServerError(err.Error(), 0)
	}

	return c.SendResponseMessage(rm, false)
}

// ok tells the client about its (client's) success.
// Returns a detailed error.
func (srv *Server) ok(c *connection.Connection) (cerr *ce.CommonError) {
	rm, err := response.NewResponse_OK()
	if err != nil {
		return ce.NewServerError(err.Error(), 0)
	}

	return c.SendResponseMessage(rm, false)
}

// closingConnection tells the client that server is going to close the
// connection.
// Returns a detailed error.
func (srv *Server) closingConnection(c *connection.Connection) (cerr *ce.CommonError) {
	rm, err := response.NewResponse_ClosingConnection()
	if err != nil {
		return ce.NewServerError(err.Error(), 0)
	}

	return c.SendResponseMessage(rm, false)
}

// showingText tells the client that server is showing text.
// Returns a detailed error.
func (srv *Server) showingText(c *connection.Connection, text string) (cerr *ce.CommonError) {
	rm, err := response.NewResponse_ShowingText(text)
	if err != nil {
		return ce.NewServerError(err.Error(), 0)
	}

	return c.SendResponseMessage(rm, false)
}

// showingBinary tells the client that server is showing binary data.
// Returns a detailed error.
func (srv *Server) showingBinary(c *connection.Connection, data []byte) (cerr *ce.CommonError) {
	rm, err := response.NewResponse_ShowingBinary(data)
	if err != nil {
		return ce.NewServerError(err.Error(), 0)
	}

	return c.SendResponseMessage(rm, true)
}

// textRecordExists tells the client that a text record exists.
// Returns a detailed error.
func (srv *Server) textRecordExists(c *connection.Connection) (cerr *ce.CommonError) {
	rm, err := response.NewResponse_TextRecordExists()
	if err != nil {
		return ce.NewServerError(err.Error(), 0)
	}

	return c.SendResponseMessage(rm, false)
}

// binaryRecordExists tells the client that a binary record exists.
// Returns a detailed error.
func (srv *Server) binaryRecordExists(c *connection.Connection) (cerr *ce.CommonError) {
	rm, err := response.NewResponse_BinaryRecordExists()
	if err != nil {
		return ce.NewServerError(err.Error(), 0)
	}

	return c.SendResponseMessage(rm, true)
}

// textRecordDoesNotExist tells the client that a text record does not exist.
// Returns a detailed error.
func (srv *Server) textRecordDoesNotExist(c *connection.Connection) (cerr *ce.CommonError) {
	rm, err := response.NewResponse_TextRecordDoesNotExist()
	if err != nil {
		return ce.NewServerError(err.Error(), 0)
	}

	return c.SendResponseMessage(rm, false)
}

// binaryRecordDoesNotExist tells the client that a binary record does not exist.
// Returns a detailed error.
func (srv *Server) binaryRecordDoesNotExist(c *connection.Connection) (cerr *ce.CommonError) {
	rm, err := response.NewResponse_BinaryRecordDoesNotExist()
	if err != nil {
		return ce.NewServerError(err.Error(), 0)
	}

	return c.SendResponseMessage(rm, true)
}

// textFileExists tells the client that a text file exists.
// Returns a detailed error.
func (srv *Server) textFileExists(c *connection.Connection) (cerr *ce.CommonError) {
	rm, err := response.NewResponse_TextFileExists()
	if err != nil {
		return ce.NewServerError(err.Error(), 0)
	}

	return c.SendResponseMessage(rm, false)
}

// binaryFileExists tells the client that a binary file exists.
// Returns a detailed error.
func (srv *Server) binaryFileExists(c *connection.Connection) (cerr *ce.CommonError) {
	rm, err := response.NewResponse_BinaryFileExists()
	if err != nil {
		return ce.NewServerError(err.Error(), 0)
	}

	return c.SendResponseMessage(rm, true)
}

// textFileDoesNotExist tells the client that a text file does not exist.
// Returns a detailed error.
func (srv *Server) textFileDoesNotExist(c *connection.Connection) (cerr *ce.CommonError) {
	rm, err := response.NewResponse_TextFileDoesNotExist()
	if err != nil {
		return ce.NewServerError(err.Error(), 0)
	}

	return c.SendResponseMessage(rm, false)
}

// binaryFileDoesNotExist tells the client that a binary file does not exist.
// Returns a detailed error.
func (srv *Server) binaryFileDoesNotExist(c *connection.Connection) (cerr *ce.CommonError) {
	rm, err := response.NewResponse_BinaryFileDoesNotExist()
	if err != nil {
		return ce.NewServerError(err.Error(), 0)
	}

	return c.SendResponseMessage(rm, true)
}
