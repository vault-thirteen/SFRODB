package client

import (
	"github.com/vault-thirteen/SFRODB/pkg/SFRODB/common/connection"
	ce "github.com/vault-thirteen/SFRODB/pkg/SFRODB/common/error"
	"github.com/vault-thirteen/SFRODB/pkg/SFRODB/common/request"
)

// closeConnection asks server to close the connection.
// Returns a detailed error.
func (cli *Client) closeConnection(con *connection.Connection) (cerr *ce.CommonError) {
	var rm *request.Request
	var err error
	rm, err = request.NewRequest_CloseConnection()
	if err != nil {
		return ce.NewClientError(err.Error(), 0, cli.id)
	}

	return con.SendRequestMessage(rm)
}

// showData asks server for data.
// Returns a detailed error.
func (cli *Client) showData(con *connection.Connection, uid string) (cerr *ce.CommonError) {
	var rm *request.Request
	var err error
	rm, err = request.NewRequest_ShowData(uid)
	if err != nil {
		return ce.NewClientError(err.Error(), 0, cli.id)
	}

	return con.SendRequestMessage(rm)
}

// searchRecord asks server to check existence of a record in cache.
// Returns a detailed error.
func (cli *Client) searchRecord(con *connection.Connection, uid string) (cerr *ce.CommonError) {
	var rm *request.Request
	var err error
	rm, err = request.NewRequest_SearchRecord(uid)
	if err != nil {
		return ce.NewClientError(err.Error(), 0, cli.id)
	}

	return con.SendRequestMessage(rm)
}

// searchFile asks server to check existence of a file.
// Returns a detailed error.
func (cli *Client) searchFile(con *connection.Connection, uid string) (cerr *ce.CommonError) {
	var rm *request.Request
	var err error
	rm, err = request.NewRequest_SearchFile(uid)
	if err != nil {
		return ce.NewClientError(err.Error(), 0, cli.id)
	}

	return con.SendRequestMessage(rm)
}

// forgetRecord asks server to remove a record from cache.
// Returns a detailed error.
func (cli *Client) forgetRecord(con *connection.Connection, uid string) (cerr *ce.CommonError) {
	var rm *request.Request
	var err error
	rm, err = request.NewRequest_ForgetRecord(uid)
	if err != nil {
		return ce.NewClientError(err.Error(), 0, cli.id)
	}

	return con.SendRequestMessage(rm)
}

// resetCache asks server to remove all records from cache.
// Returns a detailed error.
func (cli *Client) resetCache(con *connection.Connection) (cerr *ce.CommonError) {
	var rm *request.Request
	var err error
	rm, err = request.NewRequest_ResetCache()
	if err != nil {
		return ce.NewClientError(err.Error(), 0, cli.id)
	}

	return con.SendRequestMessage(rm)
}
