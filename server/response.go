package server

import (
	"github.com/vault-thirteen/SFRODB/common"
	"github.com/vault-thirteen/SFRODB/common/connection"
	ce "github.com/vault-thirteen/SFRODB/common/error"
)

// clientError tells the client about its (client's) error.
// Returns a detailed error.
func (srv *Server) clientError(c *connection.Connection) (err error) {
	var rm *common.Response
	rm, err = common.NewResponse_ClientError()
	if err != nil {
		return ce.NewServerError(err.Error(), 0)
	}

	err = c.SendResponseMessage(rm, false)
	if err != nil {
		return ce.NewServerError(err.Error(), 0)
	}

	return nil
}

// ok tells the client about its (client's) success.
// Returns a detailed error.
func (srv *Server) ok(c *connection.Connection) (err error) {
	var rm *common.Response
	rm, err = common.NewResponse_OK()
	if err != nil {
		return ce.NewServerError(err.Error(), 0)
	}

	err = c.SendResponseMessage(rm, false)
	if err != nil {
		return ce.NewServerError(err.Error(), 0)
	}

	return nil
}

// closingConnection tells the client that server is going to close the
// connection.
// Returns a detailed error.
func (srv *Server) closingConnection(c *connection.Connection) (err error) {
	var rm *common.Response
	rm, err = common.NewResponse_ClosingConnection()
	if err != nil {
		return ce.NewServerError(err.Error(), 0)
	}

	err = c.SendResponseMessage(rm, false)
	if err != nil {
		return ce.NewServerError(err.Error(), 0)
	}

	return nil
}

// showingText tells the client that server is showing text.
// Returns a detailed error.
func (srv *Server) showingText(c *connection.Connection, text string) (err error) {
	var rm *common.Response
	rm, err = common.NewResponse_ShowingText(text)
	if err != nil {
		return ce.NewServerError(err.Error(), 0)
	}

	err = c.SendResponseMessage(rm, false)
	if err != nil {
		return ce.NewServerError(err.Error(), 0)
	}

	return nil
}

// showingBinary tells the client that server is showing binary data.
// Returns a detailed error.
func (srv *Server) showingBinary(c *connection.Connection, data []byte) (err error) {
	var rm *common.Response
	rm, err = common.NewResponse_ShowingBinary(data)
	if err != nil {
		return ce.NewServerError(err.Error(), 0)
	}

	err = c.SendResponseMessage(rm, true)
	if err != nil {
		return ce.NewServerError(err.Error(), 0)
	}

	return nil
}

// textRecordExists tells the client that a text record exists.
// Returns a detailed error.
func (srv *Server) textRecordExists(c *connection.Connection) (err error) {
	var rm *common.Response
	rm, err = common.NewResponse_TextRecordExists()
	if err != nil {
		return ce.NewServerError(err.Error(), 0)
	}

	err = c.SendResponseMessage(rm, false)
	if err != nil {
		return ce.NewServerError(err.Error(), 0)
	}

	return nil
}

// binaryRecordExists tells the client that a binary record exists.
// Returns a detailed error.
func (srv *Server) binaryRecordExists(c *connection.Connection) (err error) {
	var rm *common.Response
	rm, err = common.NewResponse_BinaryRecordExists()
	if err != nil {
		return ce.NewServerError(err.Error(), 0)
	}

	err = c.SendResponseMessage(rm, true)
	if err != nil {
		return ce.NewServerError(err.Error(), 0)
	}

	return nil
}

// textRecordDoesNotExist tells the client that a text record does not exist.
// Returns a detailed error.
func (srv *Server) textRecordDoesNotExist(c *connection.Connection) (err error) {
	var rm *common.Response
	rm, err = common.NewResponse_TextRecordDoesNotExist()
	if err != nil {
		return ce.NewServerError(err.Error(), 0)
	}

	err = c.SendResponseMessage(rm, false)
	if err != nil {
		return ce.NewServerError(err.Error(), 0)
	}

	return nil
}

// binaryRecordDoesNotExist tells the client that a binary record does not exist.
// Returns a detailed error.
func (srv *Server) binaryRecordDoesNotExist(c *connection.Connection) (err error) {
	var rm *common.Response
	rm, err = common.NewResponse_BinaryRecordDoesNotExist()
	if err != nil {
		return ce.NewServerError(err.Error(), 0)
	}

	err = c.SendResponseMessage(rm, true)
	if err != nil {
		return ce.NewServerError(err.Error(), 0)
	}

	return nil
}

// textFileExists tells the client that a text file exists.
// Returns a detailed error.
func (srv *Server) textFileExists(c *connection.Connection) (err error) {
	var rm *common.Response
	rm, err = common.NewResponse_TextFileExists()
	if err != nil {
		return ce.NewServerError(err.Error(), 0)
	}

	err = c.SendResponseMessage(rm, false)
	if err != nil {
		return ce.NewServerError(err.Error(), 0)
	}

	return nil
}

// binaryFileExists tells the client that a binary file exists.
// Returns a detailed error.
func (srv *Server) binaryFileExists(c *connection.Connection) (err error) {
	var rm *common.Response
	rm, err = common.NewResponse_BinaryFileExists()
	if err != nil {
		return ce.NewServerError(err.Error(), 0)
	}

	err = c.SendResponseMessage(rm, true)
	if err != nil {
		return ce.NewServerError(err.Error(), 0)
	}

	return nil
}

// textFileDoesNotExist tells the client that a text file does not exist.
// Returns a detailed error.
func (srv *Server) textFileDoesNotExist(c *connection.Connection) (err error) {
	var rm *common.Response
	rm, err = common.NewResponse_TextFileDoesNotExist()
	if err != nil {
		return ce.NewServerError(err.Error(), 0)
	}

	err = c.SendResponseMessage(rm, false)
	if err != nil {
		return ce.NewServerError(err.Error(), 0)
	}

	return nil
}

// binaryFileDoesNotExist tells the client that a binary file does not exist.
// Returns a detailed error.
func (srv *Server) binaryFileDoesNotExist(c *connection.Connection) (err error) {
	var rm *common.Response
	rm, err = common.NewResponse_BinaryFileDoesNotExist()
	if err != nil {
		return ce.NewServerError(err.Error(), 0)
	}

	err = c.SendResponseMessage(rm, true)
	if err != nil {
		return ce.NewServerError(err.Error(), 0)
	}

	return nil
}
