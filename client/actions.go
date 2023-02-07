package client

import (
	"github.com/vault-thirteen/SFRODB/common"
	"github.com/vault-thirteen/SFRODB/common/method"
)

// CloseConnection_Main tells the server to close the main connection.
// Returns a detailed error.
func (cli *Client) CloseConnection_Main(normalExit bool) (err error) {
	return cli.closeConnection_any(true, normalExit, false)
}

// CloseConnection_Aux tells the server to close the auxiliary connection.
// Returns a detailed error.
func (cli *Client) CloseConnection_Aux(normalExit bool) (err error) {
	return cli.closeConnection_any(false, normalExit, false)
}

// closeConnection_any tells the server to close the connection.
// Returns a detailed error.
func (cli *Client) closeConnection_any(useMainConnection bool, normalExit bool, useBinary bool) (err error) {
	if useMainConnection {
		err = cli.closeConnection(cli.mainConnection)
	} else {
		err = cli.closeConnection(cli.auxConnection)
	}
	if err != nil {
		return err
	}

	// If we are closing connection due to an error, we do not wait for the
	// server's response.
	if !normalExit {
		return nil
	}

	var resp *common.Response
	if useMainConnection {
		resp, err = cli.mainConnection.GetResponseMessage(useBinary)
	} else {
		resp, err = cli.auxConnection.GetResponseMessage(useBinary)
	}
	if err != nil {
		return common.NewServerError(err.Error(), 0)
	}

	if resp.Method == method.ClosingConnection {
		return nil
	}

	return common.NewServerError(common.ErrSomethingWentWrong, resp.Method)
}

// ShowText requests a text from server and returns it.
// Returns a detailed error.
func (cli *Client) ShowText(uid string) (text string, err error) {
	err = cli.showText(cli.mainConnection, uid)
	if err != nil {
		return "", err
	}

	var resp *common.Response
	resp, err = cli.mainConnection.GetResponseMessage(false)
	if err != nil {
		return "", common.NewServerError(err.Error(), 0)
	}

	switch resp.Method {
	case method.ClientError:
		return "", common.NewClientError(common.ErrSomethingWentWrong, resp.Method)
	case method.ClosingConnection:
		return "", common.NewServerError(common.ErrSomethingWentWrong, resp.Method)
	}

	return resp.Text, nil
}

// ShowBinary requests a binary data from server and returns it.
// Returns a detailed error.
func (cli *Client) ShowBinary(uid string) (data []byte, err error) {
	err = cli.showBinary(cli.mainConnection, uid)
	if err != nil {
		return nil, err
	}

	var resp *common.Response
	resp, err = cli.mainConnection.GetResponseMessage(true)
	if err != nil {
		return nil, common.NewServerError(err.Error(), 0)
	}

	switch resp.Method {
	case method.ClientError:
		return nil, common.NewClientError(common.ErrSomethingWentWrong, resp.Method)
	case method.ClosingConnection:
		return nil, common.NewServerError(common.ErrSomethingWentWrong, resp.Method)
	}

	return resp.Data, nil
}

// SearchTextRecord asks server to check existence of a text record in cache.
// Returns a detailed error.
func (cli *Client) SearchTextRecord(uid string) (recExists bool, err error) {
	err = cli.searchTextRecord(cli.mainConnection, uid)
	if err != nil {
		return false, err
	}

	var resp *common.Response
	resp, err = cli.mainConnection.GetResponseMessage(false)
	if err != nil {
		return false, common.NewServerError(err.Error(), 0)
	}

	switch resp.Method {
	case method.TextRecordExists:
		return true, nil
	case method.TextRecordDoesNotExist:
		return false, nil
	case method.ClientError:
		return false, common.NewClientError(common.ErrSomethingWentWrong, resp.Method)
	case method.ClosingConnection:
		return false, common.NewServerError(common.ErrSomethingWentWrong, resp.Method)
	default:
		return false, common.NewServerError(common.ErrSomethingWentWrong, resp.Method)
	}
}

// SearchBinaryRecord asks server to check existence of a binary record in cache.
// Returns a detailed error.
func (cli *Client) SearchBinaryRecord(uid string) (recExists bool, err error) {
	err = cli.searchBinaryRecord(cli.mainConnection, uid)
	if err != nil {
		return false, err
	}

	var resp *common.Response
	resp, err = cli.mainConnection.GetResponseMessage(true)
	if err != nil {
		return false, common.NewServerError(err.Error(), 0)
	}

	switch resp.Method {
	case method.BinaryRecordExists:
		return true, nil
	case method.BinaryRecordDoesNotExist:
		return false, nil
	case method.ClientError:
		return false, common.NewClientError(common.ErrSomethingWentWrong, resp.Method)
	case method.ClosingConnection:
		return false, common.NewServerError(common.ErrSomethingWentWrong, resp.Method)
	default:
		return false, common.NewServerError(common.ErrSomethingWentWrong, resp.Method)
	}
}

// SearchTextFile asks server to check existence of a text file.
// Returns a detailed error.
func (cli *Client) SearchTextFile(uid string) (fileExists bool, err error) {
	err = cli.searchTextFile(cli.mainConnection, uid)
	if err != nil {
		return false, err
	}

	var resp *common.Response
	resp, err = cli.mainConnection.GetResponseMessage(false)
	if err != nil {
		return false, common.NewServerError(err.Error(), 0)
	}

	switch resp.Method {
	case method.TextFileExists:
		return true, nil
	case method.TextFileDoesNotExist:
		return false, nil
	case method.ClientError:
		return false, common.NewClientError(common.ErrSomethingWentWrong, resp.Method)
	case method.ClosingConnection:
		return false, common.NewServerError(common.ErrSomethingWentWrong, resp.Method)
	default:
		return false, common.NewServerError(common.ErrSomethingWentWrong, resp.Method)
	}
}

// SearchBinaryFile asks server to check existence of a binary file.
// Returns a detailed error.
func (cli *Client) SearchBinaryFile(uid string) (fileExists bool, err error) {
	err = cli.searchBinaryFile(cli.mainConnection, uid)
	if err != nil {
		return false, err
	}

	var resp *common.Response
	resp, err = cli.mainConnection.GetResponseMessage(true)
	if err != nil {
		return false, common.NewServerError(err.Error(), 0)
	}

	switch resp.Method {
	case method.BinaryFileExists:
		return true, nil
	case method.BinaryFileDoesNotExist:
		return false, nil
	case method.ClientError:
		return false, common.NewClientError(common.ErrSomethingWentWrong, resp.Method)
	case method.ClosingConnection:
		return false, common.NewServerError(common.ErrSomethingWentWrong, resp.Method)
	default:
		return false, common.NewServerError(common.ErrSomethingWentWrong, resp.Method)
	}
}

// ForgetTextRecord requests the server to remove a text entry from cache.
// Returns a detailed error.
func (cli *Client) ForgetTextRecord(uid string) (err error) {
	err = cli.forgetTextRecord(cli.auxConnection, uid)
	if err != nil {
		return err
	}

	var resp *common.Response
	resp, err = cli.auxConnection.GetResponseMessage(false)
	if err != nil {
		return common.NewServerError(err.Error(), 0)
	}

	if resp.Method == method.OK {
		return nil
	}

	switch resp.Method {
	case method.ClientError:
		return common.NewClientError(common.ErrSomethingWentWrong, resp.Method)
	case method.ClosingConnection:
		return common.NewServerError(common.ErrSomethingWentWrong, resp.Method)
	default:
		return common.NewServerError(common.ErrSomethingWentWrong, resp.Method)
	}
}

// ForgetBinaryRecord requests the server to remove a binary entry from cache.
// Returns a detailed error.
func (cli *Client) ForgetBinaryRecord(uid string) (err error) {
	err = cli.forgetBinaryRecord(cli.auxConnection, uid)
	if err != nil {
		return err
	}

	var resp *common.Response
	resp, err = cli.auxConnection.GetResponseMessage(true)
	if err != nil {
		return common.NewServerError(err.Error(), 0)
	}

	if resp.Method == method.OK {
		return nil
	}

	switch resp.Method {
	case method.ClientError:
		return common.NewClientError(common.ErrSomethingWentWrong, resp.Method)
	case method.ClosingConnection:
		return common.NewServerError(common.ErrSomethingWentWrong, resp.Method)
	default:
		return common.NewServerError(common.ErrSomethingWentWrong, resp.Method)
	}
}

// ResetTextCache requests the server to remove all text entries from cache.
// Returns a detailed error.
func (cli *Client) ResetTextCache() (err error) {
	err = cli.resetTextCache(cli.auxConnection)
	if err != nil {
		return err
	}

	var resp *common.Response
	resp, err = cli.auxConnection.GetResponseMessage(false)
	if err != nil {
		return common.NewServerError(err.Error(), 0)
	}

	if resp.Method == method.OK {
		return nil
	}

	switch resp.Method {
	case method.ClientError:
		return common.NewClientError(common.ErrSomethingWentWrong, resp.Method)
	case method.ClosingConnection:
		return common.NewServerError(common.ErrSomethingWentWrong, resp.Method)
	default:
		return common.NewServerError(common.ErrSomethingWentWrong, resp.Method)
	}
}

// ResetBinaryCache requests the server to remove all binary entries from cache.
// Returns a detailed error.
func (cli *Client) ResetBinaryCache() (err error) {
	err = cli.resetBinaryCache(cli.auxConnection)
	if err != nil {
		return err
	}

	var resp *common.Response
	resp, err = cli.auxConnection.GetResponseMessage(true)
	if err != nil {
		return common.NewServerError(err.Error(), 0)
	}

	if resp.Method == method.OK {
		return nil
	}

	switch resp.Method {
	case method.ClientError:
		return common.NewClientError(common.ErrSomethingWentWrong, resp.Method)
	case method.ClosingConnection:
		return common.NewServerError(common.ErrSomethingWentWrong, resp.Method)
	default:
		return common.NewServerError(common.ErrSomethingWentWrong, resp.Method)
	}
}
