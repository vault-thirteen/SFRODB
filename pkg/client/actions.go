package client

import (
	"github.com/vault-thirteen/SFRODB/pkg/common/error"
	"github.com/vault-thirteen/SFRODB/pkg/common/method"
	"github.com/vault-thirteen/SFRODB/pkg/common/response"
)

// CloseConnection_Main tells the server to close the main connection.
// Returns a detailed error.
func (cli *Client) CloseConnection_Main(normalExit bool) (cerr *ce.CommonError) {
	return cli.closeConnection_any(true, normalExit, false)
}

// CloseConnection_Aux tells the server to close the auxiliary connection.
// Returns a detailed error.
func (cli *Client) CloseConnection_Aux(normalExit bool) (cerr *ce.CommonError) {
	return cli.closeConnection_any(false, normalExit, false)
}

// closeConnection_any tells the server to close the connection.
// Returns a detailed error.
func (cli *Client) closeConnection_any(useMainConnection bool, normalExit bool, useBinary bool) (cerr *ce.CommonError) {
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
		resp, cerr = cli.mainConnection.GetResponseMessage(useBinary)
	} else {
		resp, cerr = cli.auxConnection.GetResponseMessage(useBinary)
	}
	if cerr != nil {
		return cerr
	}

	if resp.Method == method.ClosingConnection {
		return nil
	}

	return ce.NewServerError(ce.ErrSomethingWentWrong, resp.Method, cli.id)
}

// ShowText requests a text from server and returns it.
// Returns a detailed error.
func (cli *Client) ShowText(uid string) (text string, cerr *ce.CommonError) {
	cerr = cli.showText(cli.mainConnection, uid)
	if cerr != nil {
		return "", cerr
	}

	var resp *response.Response
	resp, cerr = cli.mainConnection.GetResponseMessage(false)
	if cerr != nil {
		return "", cerr
	}

	switch resp.Method {
	case method.ClientError:
		return "", ce.NewClientError(ce.ErrSomethingWentWrong, resp.Method, cli.id)
	case method.ClosingConnection:
		return "", ce.NewServerError(ce.ErrSomethingWentWrong, resp.Method, cli.id)
	}

	return resp.Text, nil
}

// ShowBinary requests a binary data from server and returns it.
// Returns a detailed error.
func (cli *Client) ShowBinary(uid string) (data []byte, cerr *ce.CommonError) {
	cerr = cli.showBinary(cli.mainConnection, uid)
	if cerr != nil {
		return nil, cerr
	}

	var resp *response.Response
	resp, cerr = cli.mainConnection.GetResponseMessage(true)
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

// SearchTextRecord asks server to check existence of a text record in cache.
// Returns a detailed error.
func (cli *Client) SearchTextRecord(uid string) (recExists bool, cerr *ce.CommonError) {
	cerr = cli.searchTextRecord(cli.mainConnection, uid)
	if cerr != nil {
		return false, cerr
	}

	var resp *response.Response
	resp, cerr = cli.mainConnection.GetResponseMessage(false)
	if cerr != nil {
		return false, cerr
	}

	switch resp.Method {
	case method.TextRecordExists:
		return true, nil
	case method.TextRecordDoesNotExist:
		return false, nil
	case method.ClientError:
		return false, ce.NewClientError(ce.ErrSomethingWentWrong, resp.Method, cli.id)
	case method.ClosingConnection:
		return false, ce.NewServerError(ce.ErrSomethingWentWrong, resp.Method, cli.id)
	default:
		return false, ce.NewServerError(ce.ErrSomethingWentWrong, resp.Method, cli.id)
	}
}

// SearchBinaryRecord asks server to check existence of a binary record in cache.
// Returns a detailed error.
func (cli *Client) SearchBinaryRecord(uid string) (recExists bool, cerr *ce.CommonError) {
	cerr = cli.searchBinaryRecord(cli.mainConnection, uid)
	if cerr != nil {
		return false, cerr
	}

	var resp *response.Response
	resp, cerr = cli.mainConnection.GetResponseMessage(true)
	if cerr != nil {
		return false, cerr
	}

	switch resp.Method {
	case method.BinaryRecordExists:
		return true, nil
	case method.BinaryRecordDoesNotExist:
		return false, nil
	case method.ClientError:
		return false, ce.NewClientError(ce.ErrSomethingWentWrong, resp.Method, cli.id)
	case method.ClosingConnection:
		return false, ce.NewServerError(ce.ErrSomethingWentWrong, resp.Method, cli.id)
	default:
		return false, ce.NewServerError(ce.ErrSomethingWentWrong, resp.Method, cli.id)
	}
}

// SearchTextFile asks server to check existence of a text file.
// Returns a detailed error.
func (cli *Client) SearchTextFile(uid string) (fileExists bool, cerr *ce.CommonError) {
	cerr = cli.searchTextFile(cli.mainConnection, uid)
	if cerr != nil {
		return false, cerr
	}

	var resp *response.Response
	resp, cerr = cli.mainConnection.GetResponseMessage(false)
	if cerr != nil {
		return false, cerr
	}

	switch resp.Method {
	case method.TextFileExists:
		return true, nil
	case method.TextFileDoesNotExist:
		return false, nil
	case method.ClientError:
		return false, ce.NewClientError(ce.ErrSomethingWentWrong, resp.Method, cli.id)
	case method.ClosingConnection:
		return false, ce.NewServerError(ce.ErrSomethingWentWrong, resp.Method, cli.id)
	default:
		return false, ce.NewServerError(ce.ErrSomethingWentWrong, resp.Method, cli.id)
	}
}

// SearchBinaryFile asks server to check existence of a binary file.
// Returns a detailed error.
func (cli *Client) SearchBinaryFile(uid string) (fileExists bool, cerr *ce.CommonError) {
	cerr = cli.searchBinaryFile(cli.mainConnection, uid)
	if cerr != nil {
		return false, cerr
	}

	var resp *response.Response
	resp, cerr = cli.mainConnection.GetResponseMessage(true)
	if cerr != nil {
		return false, cerr
	}

	switch resp.Method {
	case method.BinaryFileExists:
		return true, nil
	case method.BinaryFileDoesNotExist:
		return false, nil
	case method.ClientError:
		return false, ce.NewClientError(ce.ErrSomethingWentWrong, resp.Method, cli.id)
	case method.ClosingConnection:
		return false, ce.NewServerError(ce.ErrSomethingWentWrong, resp.Method, cli.id)
	default:
		return false, ce.NewServerError(ce.ErrSomethingWentWrong, resp.Method, cli.id)
	}
}

// ForgetTextRecord requests the server to remove a text entry from cache.
// Returns a detailed error.
func (cli *Client) ForgetTextRecord(uid string) (cerr *ce.CommonError) {
	cerr = cli.forgetTextRecord(cli.auxConnection, uid)
	if cerr != nil {
		return cerr
	}

	var resp *response.Response
	resp, cerr = cli.auxConnection.GetResponseMessage(false)
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

// ForgetBinaryRecord requests the server to remove a binary entry from cache.
// Returns a detailed error.
func (cli *Client) ForgetBinaryRecord(uid string) (cerr *ce.CommonError) {
	cerr = cli.forgetBinaryRecord(cli.auxConnection, uid)
	if cerr != nil {
		return cerr
	}

	var resp *response.Response
	resp, cerr = cli.auxConnection.GetResponseMessage(true)
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

// ResetTextCache requests the server to remove all text entries from cache.
// Returns a detailed error.
func (cli *Client) ResetTextCache() (cerr *ce.CommonError) {
	cerr = cli.resetTextCache(cli.auxConnection)
	if cerr != nil {
		return cerr
	}

	var resp *response.Response
	resp, cerr = cli.auxConnection.GetResponseMessage(false)
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

// ResetBinaryCache requests the server to remove all binary entries from cache.
// Returns a detailed error.
func (cli *Client) ResetBinaryCache() (cerr *ce.CommonError) {
	cerr = cli.resetBinaryCache(cli.auxConnection)
	if cerr != nil {
		return cerr
	}

	var resp *response.Response
	resp, cerr = cli.auxConnection.GetResponseMessage(true)
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
