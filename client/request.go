package client

import (
	"github.com/vault-thirteen/SFRODB/common/connection"
	ce "github.com/vault-thirteen/SFRODB/common/error"
	"github.com/vault-thirteen/SFRODB/common/request"
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

// showText asks server for text.
// Returns a detailed error.
func (cli *Client) showText(con *connection.Connection, uid string) (cerr *ce.CommonError) {
	var rm *request.Request
	var err error
	rm, err = request.NewRequest_ShowText(uid)
	if err != nil {
		return ce.NewClientError(err.Error(), 0, cli.id)
	}

	return con.SendRequestMessage(rm)
}

// showBinary asks server for binary data.
// Returns a detailed error.
func (cli *Client) showBinary(con *connection.Connection, uid string) (cerr *ce.CommonError) {
	var rm *request.Request
	var err error
	rm, err = request.NewRequest_ShowBinary(uid)
	if err != nil {
		return ce.NewClientError(err.Error(), 0, cli.id)
	}

	return con.SendRequestMessage(rm)
}

// searchTextRecord asks server to check existence of a text record in cache.
// Returns a detailed error.
func (cli *Client) searchTextRecord(con *connection.Connection, uid string) (cerr *ce.CommonError) {
	var rm *request.Request
	var err error
	rm, err = request.NewRequest_SearchTextRecord(uid)
	if err != nil {
		return ce.NewClientError(err.Error(), 0, cli.id)
	}

	return con.SendRequestMessage(rm)
}

// searchBinaryRecord asks server to check existence of a binary record in cache.
// Returns a detailed error.
func (cli *Client) searchBinaryRecord(con *connection.Connection, uid string) (cerr *ce.CommonError) {
	var rm *request.Request
	var err error
	rm, err = request.NewRequest_SearchBinaryRecord(uid)
	if err != nil {
		return ce.NewClientError(err.Error(), 0, cli.id)
	}

	return con.SendRequestMessage(rm)
}

// searchTextFile asks server to check existence of a text file.
// Returns a detailed error.
func (cli *Client) searchTextFile(con *connection.Connection, uid string) (cerr *ce.CommonError) {
	var rm *request.Request
	var err error
	rm, err = request.NewRequest_SearchTextFile(uid)
	if err != nil {
		return ce.NewClientError(err.Error(), 0, cli.id)
	}

	return con.SendRequestMessage(rm)
}

// searchBinaryFile asks server to check existence of a binary file.
// Returns a detailed error.
func (cli *Client) searchBinaryFile(con *connection.Connection, uid string) (cerr *ce.CommonError) {
	var rm *request.Request
	var err error
	rm, err = request.NewRequest_SearchBinaryFile(uid)
	if err != nil {
		return ce.NewClientError(err.Error(), 0, cli.id)
	}

	return con.SendRequestMessage(rm)
}

// forgetTextRecord asks server to remove a text record from cache.
// Returns a detailed error.
func (cli *Client) forgetTextRecord(con *connection.Connection, uid string) (cerr *ce.CommonError) {
	var rm *request.Request
	var err error
	rm, err = request.NewRequest_ForgetTextRecord(uid)
	if err != nil {
		return ce.NewClientError(err.Error(), 0, cli.id)
	}

	return con.SendRequestMessage(rm)
}

// forgetBinaryRecord asks server to remove a binary record from cache.
// Returns a detailed error.
func (cli *Client) forgetBinaryRecord(con *connection.Connection, uid string) (cerr *ce.CommonError) {
	var rm *request.Request
	var err error
	rm, err = request.NewRequest_ForgetBinaryRecord(uid)
	if err != nil {
		return ce.NewClientError(err.Error(), 0, cli.id)
	}

	return con.SendRequestMessage(rm)
}

// resetTextCache asks server to remove all text records from cache.
// Returns a detailed error.
func (cli *Client) resetTextCache(con *connection.Connection) (cerr *ce.CommonError) {
	var rm *request.Request
	var err error
	rm, err = request.NewRequest_ResetTextCache()
	if err != nil {
		return ce.NewClientError(err.Error(), 0, cli.id)
	}

	return con.SendRequestMessage(rm)
}

// resetBinaryCache asks server to remove all binary records from cache.
// Returns a detailed error.
func (cli *Client) resetBinaryCache(con *connection.Connection) (cerr *ce.CommonError) {
	var rm *request.Request
	var err error
	rm, err = request.NewRequest_ResetBinaryCache()
	if err != nil {
		return ce.NewClientError(err.Error(), 0, cli.id)
	}

	return con.SendRequestMessage(rm)
}
