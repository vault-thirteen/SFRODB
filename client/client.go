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
