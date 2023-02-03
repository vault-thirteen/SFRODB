package client

import (
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/vault-thirteen/SFRODB/client/settings"
	"github.com/vault-thirteen/SFRODB/common"
)

type Client struct {
	settings          *settings.Settings
	dsn               string
	addr              *net.TCPAddr
	methodNameBuffers map[common.Method][]byte
	methodValues      map[string]common.Method
	connection        *common.Connection
}

func NewClient(stn *settings.Settings) (cli *Client, err error) {
	err = stn.Check()
	if err != nil {
		return nil, err
	}

	dsn := fmt.Sprintf("%s:%d", stn.ClientHost, stn.ClientPort)

	cli = &Client{
		settings: stn,
		dsn:      dsn,
	}

	cli.addr, err = net.ResolveTCPAddr(common.LowLevelProtocol, dsn)
	if err != nil {
		return nil, err
	}

	cli.methodNameBuffers, cli.methodValues = common.InitMethods()

	return cli, nil
}

func (cli *Client) GetDsn() (dsn string) {
	return cli.dsn
}

func (cli *Client) Start() (err error) {
	var conn net.Conn
	conn, err = net.DialTCP(common.LowLevelProtocol, nil, cli.addr)
	if err != nil {
		return err
	}

	cli.connection, err = common.NewConnection(
		conn,
		&cli.methodNameBuffers,
		&cli.methodValues,
		cli.settings.ResponseMessageLengthLimit,
	)
	if err != nil {
		log.Println(err)
		return
	}

	return nil
}

func (cli *Client) Stop() (err error) {
	err = cli.connection.Break()
	if err != nil {
		return err
	}

	return nil
}

func (cli *Client) GetText(uid string) (text string, err error) {
	var rm *common.Request
	rm, err = common.NewRequest_ShowText(uid)
	if err != nil {
		return "", err
	}

	err = cli.connection.SendRequestMessage(rm)
	if err != nil {
		return "", err
	}

	var resp *common.Response
	resp, err = cli.connection.GetResponseMessage()
	if err != nil {
		return "", err
	}

	// If something goes wrong, server warns about closing the connection.
	if resp.Method == common.MethodClosingConnection {
		return "", errors.New(common.ErrSomethingWentWrong)
	}

	return resp.Text, nil
}

func (cli *Client) GetBinary(uid string) (data []byte, err error) {
	var rm *common.Request
	rm, err = common.NewRequest_ShowBinary(uid)
	if err != nil {
		return nil, err
	}

	err = cli.connection.SendRequestMessage(rm)
	if err != nil {
		return nil, err
	}

	var resp *common.Response
	resp, err = cli.connection.GetResponseMessage()
	if err != nil {
		return nil, err
	}

	// If something goes wrong, server warns about closing the connection.
	if resp.Method == common.MethodClosingConnection {
		return nil, errors.New(common.ErrSomethingWentWrong)
	}

	return resp.Data, nil
}

func (cli *Client) SayGoodbye(normalExit bool) (err error) {
	var rm *common.Request
	rm, err = common.NewRequest_CloseConnection()
	if err != nil {
		return err
	}

	err = cli.connection.SendRequestMessage(rm)
	if err != nil {
		return err
	}

	// If we are closing connection due to an error, we do not wait for the
	// server's response.
	if !normalExit {
		return nil
	}

	var resp *common.Response
	resp, err = cli.connection.GetResponseMessage()
	if err != nil {
		return err
	}

	if resp.Method != common.MethodClosingConnection {
		return errors.New(common.ErrSomethingWentWrong)
	}

	return nil
}
