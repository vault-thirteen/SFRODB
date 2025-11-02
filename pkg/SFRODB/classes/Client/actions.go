package client

import (
	ce "github.com/vault-thirteen/SFRODB/pkg/SFRODB/classes/CommonError"
	"github.com/vault-thirteen/SFRODB/pkg/SFRODB/classes/Response"
	"github.com/vault-thirteen/SFRODB/pkg/SFRODB/classes/Status"
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
		cerr = cli.request_closeConnection(cli.mainConnection)
	} else {
		cerr = cli.request_closeConnection(cli.auxConnection)
	}
	if cerr != nil {
		return cerr
	}

	// If we are closing connection due to an error, we do not wait for the
	// server's response.
	if !normalExit {
		return nil
	}

	// Wait for server's response.
	var resp *response.Response
	if useMainConnection {
		resp, cerr = cli.mainConnection.GetResponseMessage()
	} else {
		resp, cerr = cli.auxConnection.GetResponseMessage()
	}
	if cerr != nil {
		return cerr
	}

	if resp.Status != status.Status_ClosingConnection {
		return ce.NewClientError(ErrUnexpectedServerBehaviour, 0, resp.Status, cli.id)
	}

	return nil
}

// ShowData requests a data record from server and returns it.
// Returns a detailed error.
func (cli *Client) ShowData(uid string) (data []byte, cerr *ce.CommonError) {
	cerr = cli.request_showData(cli.mainConnection, uid)
	if cerr != nil {
		return nil, cerr
	}

	var resp *response.Response
	resp, cerr = cli.mainConnection.GetResponseMessage()
	if cerr != nil {
		return nil, cerr
	}

	if resp.Status != status.Status_ShowingData {
		if resp.Status == status.Status_ClientError {
			return nil, ce.NewClientError(ErrClientError, 0, resp.Status, cli.id)
		}

		return nil, ce.NewClientError(ErrUnexpectedServerBehaviour, 0, resp.Status, cli.id)
	}

	return resp.Data, nil
}

// SearchRecord asks server to check existence of a data record in cache.
// Returns a detailed error.
func (cli *Client) SearchRecord(uid string) (recExists bool, cerr *ce.CommonError) {
	cerr = cli.request_searchRecord(cli.mainConnection, uid)
	if cerr != nil {
		return false, cerr
	}

	var resp *response.Response
	resp, cerr = cli.mainConnection.GetResponseMessage()
	if cerr != nil {
		return false, cerr
	}

	switch resp.Status {
	case status.Status_RecordExists:
		return true, nil

	case status.Status_RecordDoesNotExist:
		return false, nil

	default:
		return false, ce.NewClientError(ErrUnexpectedServerBehaviour, 0, resp.Status, cli.id)
	}
}

// SearchFile asks server to check existence of a file.
// Returns a detailed error.
func (cli *Client) SearchFile(uid string) (fileExists bool, cerr *ce.CommonError) {
	cerr = cli.request_searchFile(cli.mainConnection, uid)
	if cerr != nil {
		return false, cerr
	}

	var resp *response.Response
	resp, cerr = cli.mainConnection.GetResponseMessage()
	if cerr != nil {
		return false, cerr
	}

	switch resp.Status {
	case status.Status_FileExists:
		return true, nil

	case status.Status_FileDoesNotExist:
		return false, nil

	default:
		return false, ce.NewClientError(ErrUnexpectedServerBehaviour, 0, resp.Status, cli.id)
	}
}

// ForgetRecord requests the server to remove a data entry from cache.
// Returns a detailed error.
func (cli *Client) ForgetRecord(uid string) (cerr *ce.CommonError) {
	cerr = cli.request_forgetRecord(cli.auxConnection, uid)
	if cerr != nil {
		return cerr
	}

	var resp *response.Response
	resp, cerr = cli.auxConnection.GetResponseMessage()
	if cerr != nil {
		return cerr
	}

	if resp.Status != status.Status_OK {
		return ce.NewClientError(ErrUnexpectedServerBehaviour, 0, resp.Status, cli.id)
	}

	return nil
}

// ResetCache requests the server to remove all entries from cache.
// Returns a detailed error.
func (cli *Client) ResetCache() (cerr *ce.CommonError) {
	cerr = cli.request_resetCache(cli.auxConnection)
	if cerr != nil {
		return cerr
	}

	var resp *response.Response
	resp, cerr = cli.auxConnection.GetResponseMessage()
	if cerr != nil {
		return cerr
	}

	if resp.Status != status.Status_OK {
		return ce.NewClientError(ErrUnexpectedServerBehaviour, 0, resp.Status, cli.id)
	}

	return nil
}
