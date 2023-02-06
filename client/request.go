package client

import "github.com/vault-thirteen/SFRODB/common"

// closeConnection asks server to close the connection.
// Returns a detailed error.
func (cli *Client) closeConnection(c *common.Connection) (err error) {
	var rm *common.Request
	rm, err = common.NewRequest_CloseConnection()
	if err != nil {
		return common.NewClientError(err.Error(), 0)
	}

	err = c.SendRequestMessage(rm)
	if err != nil {
		return common.NewServerError(err.Error(), 0)
	}

	return nil
}

// showText asks server for text.
// Returns a detailed error.
func (cli *Client) showText(c *common.Connection, uid string) (err error) {
	var rm *common.Request
	rm, err = common.NewRequest_ShowText(uid)
	if err != nil {
		return common.NewClientError(err.Error(), 0)
	}

	err = c.SendRequestMessage(rm)
	if err != nil {
		return common.NewServerError(err.Error(), 0)
	}

	return nil
}

// showBinary asks server for binary data.
// Returns a detailed error.
func (cli *Client) showBinary(c *common.Connection, uid string) (err error) {
	var rm *common.Request
	rm, err = common.NewRequest_ShowBinary(uid)
	if err != nil {
		return common.NewClientError(err.Error(), 0)
	}

	err = c.SendRequestMessage(rm)
	if err != nil {
		return common.NewServerError(err.Error(), 0)
	}

	return nil
}

// searchTextRecord asks server to check existence of a text record in cache.
// Returns a detailed error.
func (cli *Client) searchTextRecord(c *common.Connection, uid string) (err error) {
	var rm *common.Request
	rm, err = common.NewRequest_SearchTextRecord(uid)
	if err != nil {
		return common.NewClientError(err.Error(), 0)
	}

	err = c.SendRequestMessage(rm)
	if err != nil {
		return common.NewServerError(err.Error(), 0)
	}

	return nil
}

// searchBinaryRecord asks server to check existence of a binary record in cache.
// Returns a detailed error.
func (cli *Client) searchBinaryRecord(c *common.Connection, uid string) (err error) {
	var rm *common.Request
	rm, err = common.NewRequest_SearchBinaryRecord(uid)
	if err != nil {
		return common.NewClientError(err.Error(), 0)
	}

	err = c.SendRequestMessage(rm)
	if err != nil {
		return common.NewServerError(err.Error(), 0)
	}

	return nil
}

// searchTextFile asks server to check existence of a text file.
// Returns a detailed error.
func (cli *Client) searchTextFile(c *common.Connection, uid string) (err error) {
	var rm *common.Request
	rm, err = common.NewRequest_SearchTextFile(uid)
	if err != nil {
		return common.NewClientError(err.Error(), 0)
	}

	err = c.SendRequestMessage(rm)
	if err != nil {
		return common.NewServerError(err.Error(), 0)
	}

	return nil
}

// searchBinaryFile asks server to check existence of a binary file.
// Returns a detailed error.
func (cli *Client) searchBinaryFile(c *common.Connection, uid string) (err error) {
	var rm *common.Request
	rm, err = common.NewRequest_SearchBinaryFile(uid)
	if err != nil {
		return common.NewClientError(err.Error(), 0)
	}

	err = c.SendRequestMessage(rm)
	if err != nil {
		return common.NewServerError(err.Error(), 0)
	}

	return nil
}

// forgetTextRecord asks server to remove a text record from cache.
// Returns a detailed error.
func (cli *Client) forgetTextRecord(c *common.Connection, uid string) (err error) {
	var rm *common.Request
	rm, err = common.NewRequest_ForgetTextRecord(uid)
	if err != nil {
		return common.NewClientError(err.Error(), 0)
	}

	err = c.SendRequestMessage(rm)
	if err != nil {
		return common.NewServerError(err.Error(), 0)
	}

	return nil
}

// forgetBinaryRecord asks server to remove a binary record from cache.
// Returns a detailed error.
func (cli *Client) forgetBinaryRecord(c *common.Connection, uid string) (err error) {
	var rm *common.Request
	rm, err = common.NewRequest_ForgetBinaryRecord(uid)
	if err != nil {
		return common.NewClientError(err.Error(), 0)
	}

	err = c.SendRequestMessage(rm)
	if err != nil {
		return common.NewServerError(err.Error(), 0)
	}

	return nil
}

// resetTextCache asks server to remove all text records from cache.
// Returns a detailed error.
func (cli *Client) resetTextCache(c *common.Connection) (err error) {
	var rm *common.Request
	rm, err = common.NewRequest_ResetTextCache()
	if err != nil {
		return common.NewClientError(err.Error(), 0)
	}

	err = c.SendRequestMessage(rm)
	if err != nil {
		return common.NewServerError(err.Error(), 0)
	}

	return nil
}

// resetBinaryCache asks server to remove all binary records from cache.
// Returns a detailed error.
func (cli *Client) resetBinaryCache(c *common.Connection) (err error) {
	var rm *common.Request
	rm, err = common.NewRequest_ResetBinaryCache()
	if err != nil {
		return common.NewClientError(err.Error(), 0)
	}

	err = c.SendRequestMessage(rm)
	if err != nil {
		return common.NewServerError(err.Error(), 0)
	}

	return nil
}
