package server

import (
	ce "github.com/vault-thirteen/SFRODB/pkg/SFRODB/classes/CommonError"
	"github.com/vault-thirteen/SFRODB/pkg/SFRODB/classes/Connection"
	"github.com/vault-thirteen/SFRODB/pkg/SFRODB/classes/Response"
)

// respond_clientError tells the client about its (client's) error.
// Returns a detailed error.
func (srv *Server) respond_clientError(con *connection.Connection) (cerr *ce.CommonError) {
	resp, err := response.New_ClientError()
	if err != nil {
		return ce.NewServerError(err.Error(), 0, 0, con.ClientId())
	}

	return con.SendResponseMessage(resp)
}

// respond_ok tells the client about its (client's) success.
// Returns a detailed error.
func (srv *Server) respond_ok(con *connection.Connection) (cerr *ce.CommonError) {
	resp, err := response.New_OK()
	if err != nil {
		return ce.NewServerError(err.Error(), 0, 0, con.ClientId())
	}

	return con.SendResponseMessage(resp)
}

// respond_closingConnection tells the client that server is going to close the
// connection.
// Returns a detailed error.
func (srv *Server) respond_closingConnection(con *connection.Connection) (cerr *ce.CommonError) {
	resp, err := response.New_ClosingConnection()
	if err != nil {
		return ce.NewServerError(err.Error(), 0, 0, con.ClientId())
	}

	return con.SendResponseMessage(resp)
}

// respond_showingData tells the client that server is showing data.
// Returns a detailed error.
func (srv *Server) respond_showingData(con *connection.Connection, data []byte) (cerr *ce.CommonError) {
	resp, err := response.New_ShowingData(data)
	if err != nil {
		return ce.NewServerError(err.Error(), 0, 0, con.ClientId())
	}

	return con.SendResponseMessage(resp)
}

// respond_recordExists tells the client that a record exists.
// Returns a detailed error.
func (srv *Server) respond_recordExists(con *connection.Connection) (cerr *ce.CommonError) {
	resp, err := response.New_RecordExists()
	if err != nil {
		return ce.NewServerError(err.Error(), 0, 0, con.ClientId())
	}

	return con.SendResponseMessage(resp)
}

// respond_recordDoesNotExist tells the client that a record does not exist.
// Returns a detailed error.
func (srv *Server) respond_recordDoesNotExist(con *connection.Connection) (cerr *ce.CommonError) {
	resp, err := response.New_RecordDoesNotExist()
	if err != nil {
		return ce.NewServerError(err.Error(), 0, 0, con.ClientId())
	}

	return con.SendResponseMessage(resp)
}

// respond_fileExists tells the client that a file exists.
// Returns a detailed error.
func (srv *Server) respond_fileExists(con *connection.Connection) (cerr *ce.CommonError) {
	resp, err := response.New_FileExists()
	if err != nil {
		return ce.NewServerError(err.Error(), 0, 0, con.ClientId())
	}

	return con.SendResponseMessage(resp)
}

// respond_fileDoesNotExist tells the client that a file does not exist.
// Returns a detailed error.
func (srv *Server) respond_fileDoesNotExist(con *connection.Connection) (cerr *ce.CommonError) {
	resp, err := response.New_FileDoesNotExist()
	if err != nil {
		return ce.NewServerError(err.Error(), 0, 0, con.ClientId())
	}

	return con.SendResponseMessage(resp)
}
