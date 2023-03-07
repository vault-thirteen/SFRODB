package client

import (
	"github.com/vault-thirteen/SFRODB/pkg/common/error"
	"github.com/vault-thirteen/SFRODB/pkg/common/method"
	"github.com/vault-thirteen/SFRODB/pkg/common/response"
)

// CloseConnection_Main tells the server to close the main connection.
// Returns a detailed error.
func (cli *Client) CloseConnection_Main(normalExit bool) (cerr *ce.CommonError) {
	return cli.closeConnection_any(true, normalExit)
}

// CloseConnection_Aux tells the server to close the auxiliary connection.
// Returns a detailed error.
func (cli *Client) CloseConnection_Aux(normalExit bool) (cerr *ce.CommonError) {
	return cli.closeConnection_any(false, normalExit)
}

// closeConnection_any tells the server to close the connection.
// Returns a detailed error.
func (cli *Client) closeConnection_any(useMainConnection bool, normalExit bool) (cerr *ce.CommonError) {
	if useMainConnection {
		cerr = cli.closeConnection(cli.mainConnection)
	} else {
		cerr = cli.closeConnection(cli.auxConnection)
	}
	if cerr != nil {
		return cerr
	}

	// If we are closing connection due to an error, we do not wait for the
	// server's response.
	if !normalExit {
		return nil
	}

	var resp *response.Response
	if useMainConnection {
		resp, cerr = cli.mainConnection.GetResponseMessage()
	} else {
		resp, cerr = cli.auxConnection.GetResponseMessage()
	}
	if cerr != nil {
		return cerr
	}

	if resp.Method == method.ClosingConnection {
		return nil
	}

	return ce.NewServerError(ce.ErrSomethingWentWrong, resp.Method, cli.id)
}

// ShowData requests a data record from server and returns it.
// Returns a detailed error.
func (cli *Client) ShowData(uid string) (data []byte, cerr *ce.CommonError) {
	cerr = cli.showData(cli.mainConnection, uid)
	if cerr != nil {
		return nil, cerr
	}

	var resp *response.Response
	resp, cerr = cli.mainConnection.GetResponseMessage()
	if cerr != nil {
		return nil, cerr
	}

	switch resp.Method {
	case method.ClientError:
		return nil, ce.NewClientError(ce.ErrSomethingWentWrong, resp.Method, cli.id)
	case method.ClosingConnection:
		return nil, ce.NewServerError(ce.ErrSomethingWentWrong, resp.Method, cli.id)
	}

	return resp.Data, nil
}

// SearchRecord asks server to check existence of a data record in cache.
// Returns a detailed error.
func (cli *Client) SearchRecord(uid string) (recExists bool, cerr *ce.CommonError) {
	cerr = cli.searchRecord(cli.mainConnection, uid)
	if cerr != nil {
		return false, cerr
	}

	var resp *response.Response
	resp, cerr = cli.mainConnection.GetResponseMessage()
	if cerr != nil {
		return false, cerr
	}

	switch resp.Method {
	case method.RecordExists:
		return true, nil
	case method.RecordDoesNotExist:
		return false, nil
	case method.ClientError:
		return false, ce.NewClientError(ce.ErrSomethingWentWrong, resp.Method, cli.id)
	case method.ClosingConnection:
		return false, ce.NewServerError(ce.ErrSomethingWentWrong, resp.Method, cli.id)
	default:
		return false, ce.NewServerError(ce.ErrSomethingWentWrong, resp.Method, cli.id)
	}
}

// SearchFile asks server to check existence of a file.
// Returns a detailed error.
func (cli *Client) SearchFile(uid string) (fileExists bool, cerr *ce.CommonError) {
	cerr = cli.searchFile(cli.mainConnection, uid)
	if cerr != nil {
		return false, cerr
	}

	var resp *response.Response
	resp, cerr = cli.mainConnection.GetResponseMessage()
	if cerr != nil {
		return false, cerr
	}

	switch resp.Method {
	case method.FileExists:
		return true, nil
	case method.FileDoesNotExist:
		return false, nil
	case method.ClientError:
		return false, ce.NewClientError(ce.ErrSomethingWentWrong, resp.Method, cli.id)
	case method.ClosingConnection:
		return false, ce.NewServerError(ce.ErrSomethingWentWrong, resp.Method, cli.id)
	default:
		return false, ce.NewServerError(ce.ErrSomethingWentWrong, resp.Method, cli.id)
	}
}

// ForgetRecord requests the server to remove a data entry from cache.
// Returns a detailed error.
func (cli *Client) ForgetRecord(uid string) (cerr *ce.CommonError) {
	cerr = cli.forgetRecord(cli.auxConnection, uid)
	if cerr != nil {
		return cerr
	}

	var resp *response.Response
	resp, cerr = cli.auxConnection.GetResponseMessage()
	if cerr != nil {
		return cerr
	}

	if resp.Method == method.OK {
		return nil
	}

	switch resp.Method {
	case method.ClientError:
		return ce.NewClientError(ce.ErrSomethingWentWrong, resp.Method, cli.id)
	case method.ClosingConnection:
		return ce.NewServerError(ce.ErrSomethingWentWrong, resp.Method, cli.id)
	default:
		return ce.NewServerError(ce.ErrSomethingWentWrong, resp.Method, cli.id)
	}
}

// ResetCache requests the server to remove all entries from cache.
// Returns a detailed error.
func (cli *Client) ResetCache() (cerr *ce.CommonError) {
	cerr = cli.resetCache(cli.auxConnection)
	if cerr != nil {
		return cerr
	}

	var resp *response.Response
	resp, cerr = cli.auxConnection.GetResponseMessage()
	if cerr != nil {
		return cerr
	}

	if resp.Method == method.OK {
		return nil
	}

	switch resp.Method {
	case method.ClientError:
		return ce.NewClientError(ce.ErrSomethingWentWrong, resp.Method, cli.id)
	case method.ClosingConnection:
		return ce.NewServerError(ce.ErrSomethingWentWrong, resp.Method, cli.id)
	default:
		return ce.NewServerError(ce.ErrSomethingWentWrong, resp.Method, cli.id)
	}
}
