package client

import (
	"fmt"
	"net"

	"github.com/vault-thirteen/SFRODB/client/settings"
	"github.com/vault-thirteen/SFRODB/common/connection"
	"github.com/vault-thirteen/SFRODB/common/method"
	"github.com/vault-thirteen/SFRODB/common/protocol"
)

// Client is client.
type Client struct {
	settings *settings.Settings

	mainDsn string
	auxDsn  string

	mainAddr *net.TCPAddr
	auxAddr  *net.TCPAddr

	methodNameBuffers map[method.Method][]byte
	methodValues      map[string]method.Method

	mainConnection *connection.Connection
	auxConnection  *connection.Connection
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

	cli.mainAddr, err = net.ResolveTCPAddr(proto.LowLevelProtocol, cli.mainDsn)
	if err != nil {
		return nil, err
	}

	cli.auxAddr, err = net.ResolveTCPAddr(proto.LowLevelProtocol, cli.auxDsn)
	if err != nil {
		return nil, err
	}

	cli.methodNameBuffers, cli.methodValues = method.InitMethods()

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
	mainConn, err = net.DialTCP(proto.LowLevelProtocol, nil, cli.mainAddr)
	if err != nil {
		return err
	}

	cli.mainConnection = connection.NewConnection(
		mainConn,
		&cli.methodNameBuffers,
		&cli.methodValues,
		cli.settings.ResponseMessageLengthLimit,
	)

	var auxConn net.Conn
	auxConn, err = net.DialTCP(proto.LowLevelProtocol, nil, cli.auxAddr)
	if err != nil {
		return err
	}

	cli.auxConnection = connection.NewConnection(
		auxConn,
		&cli.methodNameBuffers,
		&cli.methodValues,
		cli.settings.ResponseMessageLengthLimit,
	)

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

// Restart re-starts the client.
func (cli *Client) Restart(forcibly bool) (err error) {
	if forcibly {
		_ = cli.Stop()
	} else {
		err = cli.Stop()
		if err != nil {
			return err
		}
	}

	return cli.Start()
}
