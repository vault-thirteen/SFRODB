package client

import (
	ce "github.com/vault-thirteen/SFRODB/pkg/SFRODB/classes/CommonError"
	"github.com/vault-thirteen/SFRODB/pkg/SFRODB/classes/Connection"
	"github.com/vault-thirteen/SFRODB/pkg/SFRODB/classes/Request"
)

// request_closeConnection asks server to close the connection.
// Returns a detailed error.
func (cli *Client) request_closeConnection(con *connection.Connection) (cerr *ce.CommonError) {
	req, err := request.New_CloseConnection()
	if err != nil {
		return ce.NewClientError(err.Error(), 0, 0, cli.id)
	}

	return con.SendRequestMessage(req)
}

// request_showData asks server for data.
// Returns a detailed error.
func (cli *Client) request_showData(con *connection.Connection, uid string) (cerr *ce.CommonError) {
	req, err := request.New_ShowData(uid)
	if err != nil {
		return ce.NewClientError(err.Error(), 0, 0, cli.id)
	}

	return con.SendRequestMessage(req)
}

// request_searchRecord asks server to check existence of a record in cache.
// Returns a detailed error.
func (cli *Client) request_searchRecord(con *connection.Connection, uid string) (cerr *ce.CommonError) {
	req, err := request.New_SearchRecord(uid)
	if err != nil {
		return ce.NewClientError(err.Error(), 0, 0, cli.id)
	}

	return con.SendRequestMessage(req)
}

// request_searchFile asks server to check existence of a file.
// Returns a detailed error.
func (cli *Client) request_searchFile(con *connection.Connection, uid string) (cerr *ce.CommonError) {
	req, err := request.New_SearchFile(uid)
	if err != nil {
		return ce.NewClientError(err.Error(), 0, 0, cli.id)
	}

	return con.SendRequestMessage(req)
}

// request_forgetRecord asks server to remove a record from cache.
// Returns a detailed error.
func (cli *Client) request_forgetRecord(con *connection.Connection, uid string) (cerr *ce.CommonError) {
	req, err := request.New_ForgetRecord(uid)
	if err != nil {
		return ce.NewClientError(err.Error(), 0, 0, cli.id)
	}

	return con.SendRequestMessage(req)
}

// request_resetCache asks server to remove all records from cache.
// Returns a detailed error.
func (cli *Client) request_resetCache(con *connection.Connection) (cerr *ce.CommonError) {
	req, err := request.New_ResetCache()
	if err != nil {
		return ce.NewClientError(err.Error(), 0, 0, cli.id)
	}

	return con.SendRequestMessage(req)
}
