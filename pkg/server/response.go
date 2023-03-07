package server

import (
	"github.com/vault-thirteen/SFRODB/pkg/common/connection"
	"github.com/vault-thirteen/SFRODB/pkg/common/error"
	"github.com/vault-thirteen/SFRODB/pkg/common/response"
)

// clientError tells the client about its (client's) error.
// Returns a detailed error.
func (srv *Server) clientError(con *connection.Connection) (cerr *ce.CommonError) {
	rm, err := response.NewResponse_ClientError()
	if err != nil {
		return ce.NewServerError(err.Error(), 0, con.GetClientId())
	}

	return con.SendResponseMessage(rm)
}

// ok tells the client about its (client's) success.
// Returns a detailed error.
func (srv *Server) ok(con *connection.Connection) (cerr *ce.CommonError) {
	rm, err := response.NewResponse_OK()
	if err != nil {
		return ce.NewServerError(err.Error(), 0, con.GetClientId())
	}

	return con.SendResponseMessage(rm)
}

// closingConnection tells the client that server is going to close the
// connection.
// Returns a detailed error.
func (srv *Server) closingConnection(con *connection.Connection) (cerr *ce.CommonError) {
	rm, err := response.NewResponse_ClosingConnection()
	if err != nil {
		return ce.NewServerError(err.Error(), 0, con.GetClientId())
	}

	return con.SendResponseMessage(rm)
}

// showingData tells the client that server is showing data.
// Returns a detailed error.
func (srv *Server) showingData(con *connection.Connection, data []byte) (cerr *ce.CommonError) {
	rm, err := response.NewResponse_ShowingData(data)
	if err != nil {
		return ce.NewServerError(err.Error(), 0, con.GetClientId())
	}

	return con.SendResponseMessage(rm)
}

// recordExists tells the client that a record exists.
// Returns a detailed error.
func (srv *Server) recordExists(con *connection.Connection) (cerr *ce.CommonError) {
	rm, err := response.NewResponse_RecordExists()
	if err != nil {
		return ce.NewServerError(err.Error(), 0, con.GetClientId())
	}

	return con.SendResponseMessage(rm)
}

// recordDoesNotExist tells the client that a record does not exist.
// Returns a detailed error.
func (srv *Server) recordDoesNotExist(con *connection.Connection) (cerr *ce.CommonError) {
	rm, err := response.NewResponse_RecordDoesNotExist()
	if err != nil {
		return ce.NewServerError(err.Error(), 0, con.GetClientId())
	}

	return con.SendResponseMessage(rm)
}

// fileExists tells the client that a file exists.
// Returns a detailed error.
func (srv *Server) fileExists(con *connection.Connection) (cerr *ce.CommonError) {
	rm, err := response.NewResponse_FileExists()
	if err != nil {
		return ce.NewServerError(err.Error(), 0, con.GetClientId())
	}

	return con.SendResponseMessage(rm)
}

// fileDoesNotExist tells the client that a file does not exist.
// Returns a detailed error.
func (srv *Server) fileDoesNotExist(con *connection.Connection) (cerr *ce.CommonError) {
	rm, err := response.NewResponse_FileDoesNotExist()
	if err != nil {
		return ce.NewServerError(err.Error(), 0, con.GetClientId())
	}

	return con.SendResponseMessage(rm)
}
