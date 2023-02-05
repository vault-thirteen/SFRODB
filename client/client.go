package client

import (
	"fmt"
	"log"
	"net"

	"github.com/vault-thirteen/SFRODB/client/settings"
	"github.com/vault-thirteen/SFRODB/common"
)

// Client is client.
type Client struct {
	settings *settings.Settings

	mainDsn string
	auxDsn  string

	mainAddr *net.TCPAddr
	auxAddr  *net.TCPAddr

	methodNameBuffers map[common.Method][]byte
	methodValues      map[string]common.Method

	mainConnection *common.Connection
	auxConnection  *common.Connection
}

// NewClient creates a client.
func NewClient(stn *settings.Settings) (cli *Client, err error) {
	err = stn.Check()
	if err != nil {
		return nil, err
	}

	cli = &Client{
		settings: stn,
		mainDsn:  fmt.Sprintf("%s:%d", stn.Host, stn.MainPort),
		auxDsn:   fmt.Sprintf("%s:%d", stn.Host, stn.AuxPort),
	}

	cli.mainAddr, err = net.ResolveTCPAddr(common.LowLevelProtocol, cli.mainDsn)
	if err != nil {
		return nil, err
	}

	cli.auxAddr, err = net.ResolveTCPAddr(common.LowLevelProtocol, cli.auxDsn)
	if err != nil {
		return nil, err
	}

	cli.methodNameBuffers, cli.methodValues = common.InitMethods()

	return cli, nil
}

// GetMainDsn returns the DSN of the main connection.
func (cli *Client) GetMainDsn() (dsn string) {
	return cli.mainDsn
}

// GetAuxDsn returns the DSN of the auxiliary connection.
func (cli *Client) GetAuxDsn() (dsn string) {
	return cli.auxDsn
}

// Start starts the client.
func (cli *Client) Start() (err error) {
	var mainConn net.Conn
	mainConn, err = net.DialTCP(common.LowLevelProtocol, nil, cli.mainAddr)
	if err != nil {
		return err
	}

	cli.mainConnection, err = common.NewConnection(
		mainConn,
		&cli.methodNameBuffers,
		&cli.methodValues,
		cli.settings.ResponseMessageLengthLimit,
	)
	if err != nil {
		log.Println(err)
		return err
	}

	var auxConn net.Conn
	auxConn, err = net.DialTCP(common.LowLevelProtocol, nil, cli.auxAddr)
	if err != nil {
		return err
	}

	cli.auxConnection, err = common.NewConnection(
		auxConn,
		&cli.methodNameBuffers,
		&cli.methodValues,
		cli.settings.ResponseMessageLengthLimit,
	)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// Stop stops the client.
func (cli *Client) Stop() (err error) {
	err = cli.mainConnection.Break()
	if err != nil {
		return err
	}

	err = cli.auxConnection.Break()
	if err != nil {
		return err
	}

	return nil
}

// GetText requests a text from server and returns it.
// Returns a detailed error.
func (cli *Client) GetText(uid string) (text string, err error) {
	var rm *common.Request
	rm, err = common.NewRequest_ShowText(uid)
	if err != nil {
		return "", common.NewClientError(err.Error(), 0)
	}

	err = cli.mainConnection.SendRequestMessage(rm)
	if err != nil {
		return "", common.NewServerError(err.Error(), 0)
	}

	var resp *common.Response
	resp, err = cli.mainConnection.GetResponseMessage()
	if err != nil {
		return "", common.NewServerError(err.Error(), 0)
	}

	switch resp.Method {
	case common.MethodClientError:
		return "", common.NewClientError(common.ErrSomethingWentWrong, resp.Method)
	case common.MethodClosingConnection:
		return "", common.NewServerError(common.ErrSomethingWentWrong, resp.Method)
	}

	return resp.Text, nil
}

// GetBinary requests a binary data from server and returns it.
// Returns a detailed error.
func (cli *Client) GetBinary(uid string) (data []byte, err error) {
	var rm *common.Request
	rm, err = common.NewRequest_ShowBinary(uid)
	if err != nil {
		return nil, common.NewClientError(err.Error(), 0)
	}

	err = cli.mainConnection.SendRequestMessage(rm)
	if err != nil {
		return nil, common.NewServerError(err.Error(), 0)
	}

	var resp *common.Response
	resp, err = cli.mainConnection.GetResponseMessage()
	if err != nil {
		return nil, common.NewServerError(err.Error(), 0)
	}

	switch resp.Method {
	case common.MethodClientError:
		return nil, common.NewClientError(common.ErrSomethingWentWrong, resp.Method)
	case common.MethodClosingConnection:
		return nil, common.NewServerError(common.ErrSomethingWentWrong, resp.Method)
	}

	return resp.Data, nil
}

// RemoveText requests the server to remove a text entry from cache.
// Returns a detailed error.
func (cli *Client) RemoveText(uid string) (err error) {
	var rm *common.Request
	rm, err = common.NewRequest_RemoveText(uid)
	if err != nil {
		return common.NewClientError(err.Error(), 0)
	}

	err = cli.auxConnection.SendRequestMessage(rm)
	if err != nil {
		return common.NewServerError(err.Error(), 0)
	}

	var resp *common.Response
	resp, err = cli.auxConnection.GetResponseMessage()
	if err != nil {
		return common.NewServerError(err.Error(), 0)
	}

	if resp.Method == common.MethodOK {
		return nil
	}

	switch resp.Method {
	case common.MethodClientError:
		return common.NewClientError(common.ErrSomethingWentWrong, resp.Method)
	case common.MethodClosingConnection:
		return common.NewServerError(common.ErrSomethingWentWrong, resp.Method)
	default:
		return common.NewServerError(common.ErrSomethingWentWrong, resp.Method)
	}
}

// RemoveBinary requests the server to remove a binary entry from cache.
// Returns a detailed error.
func (cli *Client) RemoveBinary(uid string) (err error) {
	var rm *common.Request
	rm, err = common.NewRequest_RemoveBinary(uid)
	if err != nil {
		return common.NewClientError(err.Error(), 0)
	}

	err = cli.auxConnection.SendRequestMessage(rm)
	if err != nil {
		return common.NewServerError(err.Error(), 0)
	}

	var resp *common.Response
	resp, err = cli.auxConnection.GetResponseMessage()
	if err != nil {
		return common.NewServerError(err.Error(), 0)
	}

	if resp.Method == common.MethodOK {
		return nil
	}

	switch resp.Method {
	case common.MethodClientError:
		return common.NewClientError(common.ErrSomethingWentWrong, resp.Method)
	case common.MethodClosingConnection:
		return common.NewServerError(common.ErrSomethingWentWrong, resp.Method)
	default:
		return common.NewServerError(common.ErrSomethingWentWrong, resp.Method)
	}
}

// ClearTextCache requests the server to remove all text entries from cache.
// Returns a detailed error.
func (cli *Client) ClearTextCache() (err error) {
	var rm *common.Request
	rm, err = common.NewRequest_ClearTextCache()
	if err != nil {
		return common.NewClientError(err.Error(), 0)
	}

	err = cli.auxConnection.SendRequestMessage(rm)
	if err != nil {
		return common.NewServerError(err.Error(), 0)
	}

	var resp *common.Response
	resp, err = cli.auxConnection.GetResponseMessage()
	if err != nil {
		return common.NewServerError(err.Error(), 0)
	}

	if resp.Method == common.MethodOK {
		return nil
	}

	switch resp.Method {
	case common.MethodClientError:
		return common.NewClientError(common.ErrSomethingWentWrong, resp.Method)
	case common.MethodClosingConnection:
		return common.NewServerError(common.ErrSomethingWentWrong, resp.Method)
	default:
		return common.NewServerError(common.ErrSomethingWentWrong, resp.Method)
	}
}

// ClearBinaryCache requests the server to remove all binary entries from cache.
// Returns a detailed error.
func (cli *Client) ClearBinaryCache() (err error) {
	var rm *common.Request
	rm, err = common.NewRequest_ClearBinaryCache()
	if err != nil {
		return common.NewClientError(err.Error(), 0)
	}

	err = cli.auxConnection.SendRequestMessage(rm)
	if err != nil {
		return common.NewServerError(err.Error(), 0)
	}

	var resp *common.Response
	resp, err = cli.auxConnection.GetResponseMessage()
	if err != nil {
		return common.NewServerError(err.Error(), 0)
	}

	if resp.Method == common.MethodOK {
		return nil
	}

	switch resp.Method {
	case common.MethodClientError:
		return common.NewClientError(common.ErrSomethingWentWrong, resp.Method)
	case common.MethodClosingConnection:
		return common.NewServerError(common.ErrSomethingWentWrong, resp.Method)
	default:
		return common.NewServerError(common.ErrSomethingWentWrong, resp.Method)
	}
}

// SayGoodbyeOnMain tells the server to close the main connection.
// Returns a detailed error.
func (cli *Client) SayGoodbyeOnMain(normalExit bool) (err error) {
	return cli.sayGoodbye(true, normalExit)
}

// SayGoodbyeOnAux tells the server to close the auxiliary connection.
// Returns a detailed error.
func (cli *Client) SayGoodbyeOnAux(normalExit bool) (err error) {
	return cli.sayGoodbye(false, normalExit)
}

// sayGoodbye tells the server to close the connection.
// Returns a detailed error.
func (cli *Client) sayGoodbye(useMainConnection bool, normalExit bool) (err error) {
	var rm *common.Request
	rm, err = common.NewRequest_CloseConnection()
	if err != nil {
		return common.NewClientError(err.Error(), 0)
	}

	if useMainConnection {
		err = cli.mainConnection.SendRequestMessage(rm)
	} else {
		err = cli.auxConnection.SendRequestMessage(rm)
	}
	if err != nil {
		return common.NewServerError(err.Error(), 0)
	}

	// If we are closing connection due to an error, we do not wait for the
	// server's response.
	if !normalExit {
		return nil
	}

	var resp *common.Response
	if useMainConnection {
		resp, err = cli.mainConnection.GetResponseMessage()
	} else {
		resp, err = cli.auxConnection.GetResponseMessage()
	}
	if err != nil {
		return common.NewServerError(err.Error(), 0)
	}

	if resp.Method == common.MethodClosingConnection {
		return nil
	}

	return common.NewServerError(common.ErrSomethingWentWrong, resp.Method)
}
