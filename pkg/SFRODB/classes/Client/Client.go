package client

import (
	"fmt"
	"net"
	"sync"
	"sync/atomic"

	cs "github.com/vault-thirteen/SFRODB/pkg/SFRODB/classes/ClientSettings"
	ce "github.com/vault-thirteen/SFRODB/pkg/SFRODB/classes/CommonError"
	"github.com/vault-thirteen/SFRODB/pkg/SFRODB/classes/Connection"
	"github.com/vault-thirteen/SFRODB/pkg/SFRODB/protocol"
	"github.com/vault-thirteen/SFRODB/pkg/SFRODB/std/tcp"
	ae "github.com/vault-thirteen/auxie/errors"
)

const (
	// ClientIdNone is a client ID for a non-network operations.
	ClientIdNone = ""

	// ClientIdIncoming is a client ID for a request incoming to server.
	ClientIdIncoming = "s"
)

const (
	ErrDoubleStartIsNotPossible  = "double start is not possible"
	ErrDoubleStopIsNotPossible   = "double stop is not possible"
	ErrUnexpectedServerBehaviour = "unexpected server behaviour"
	ErrClientError               = "client error"
)

// Client is client.
type Client struct {
	// All the clients must have unique IDs.
	id string

	settings *cs.ClientSettings

	mainDsn string
	auxDsn  string

	mainAddr *net.TCPAddr
	auxAddr  *net.TCPAddr

	mainConnection *connection.Connection
	auxConnection  *connection.Connection

	// Internal control structures.
	startStopLock *sync.Mutex
	isWorking     *atomic.Bool
}

func New(stn *cs.ClientSettings, id string) (cli *Client, err error) {
	err = stn.Check()
	if err != nil {
		return nil, err
	}

	cli = &Client{
		id:            id,
		settings:      stn,
		mainDsn:       fmt.Sprintf("%s:%d", stn.Host, stn.MainPort),
		auxDsn:        fmt.Sprintf("%s:%d", stn.Host, stn.AuxPort),
		startStopLock: new(sync.Mutex),
		isWorking:     new(atomic.Bool),
	}

	cli.mainAddr, err = net.ResolveTCPAddr(protocol.LowLevelProtocol, cli.mainDsn)
	if err != nil {
		return nil, err
	}

	cli.auxAddr, err = net.ResolveTCPAddr(protocol.LowLevelProtocol, cli.auxDsn)
	if err != nil {
		return nil, err
	}

	return cli, nil
}

// GetId returns the ID of the client.
func (cli *Client) GetId() (id string) {
	return cli.id
}

// GetMainDsn returns the DSN of the main connection of the client.
func (cli *Client) GetMainDsn() (dsn string) {
	return cli.mainDsn
}

// GetAuxDsn returns the DSN of the auxiliary connection of the client.
func (cli *Client) GetAuxDsn() (dsn string) {
	return cli.auxDsn
}

// Start starts the client.
func (cli *Client) Start() (cerr *ce.CommonError) {
	cli.startStopLock.Lock()
	defer cli.startStopLock.Unlock()

	return cli.start()
}

func (cli *Client) start() (cerr *ce.CommonError) {
	if cli.isWorking.Load() {
		return ce.NewClientError(ErrDoubleStartIsNotPossible, 0, 0, ClientIdNone)
	}

	cerr = cli.startMainConnection()
	if cerr != nil {
		return cerr
	}

	cerr = cli.startAuxConnection()
	if cerr != nil {
		return cerr
	}

	cli.isWorking.Store(true)

	return nil
}

func (cli *Client) startMainConnection() (cerr *ce.CommonError) {
	mainConn, err := net.DialTCP(protocol.LowLevelProtocol, nil, cli.mainAddr)
	if err != nil {
		return ce.NewClientError(err.Error(), 0, 0, cli.id)
	}

	err = tcp.EnableKeepAlives(mainConn, protocol.TcpKeepAliveIsEnabled, protocol.TcpKeepAlivePeriodSec)
	if err != nil {
		closeErr := mainConn.Close()
		if closeErr != nil {
			err = ae.Combine(err, closeErr)
		}
		return ce.NewClientError(err.Error(), 0, 0, cli.id)
	}

	cli.mainConnection = connection.New(mainConn, cli.settings.ResponseMessageLengthLimit, cli.id)

	return nil
}

func (cli *Client) startAuxConnection() (cerr *ce.CommonError) {
	auxConn, err := net.DialTCP(protocol.LowLevelProtocol, nil, cli.auxAddr)
	if err != nil {
		return ce.NewClientError(err.Error(), 0, 0, cli.id)
	}

	err = tcp.EnableKeepAlives(auxConn, protocol.TcpKeepAliveIsEnabled, protocol.TcpKeepAlivePeriodSec)
	if err != nil {
		closeErr := auxConn.Close()
		if closeErr != nil {
			err = ae.Combine(err, closeErr)
		}
		return ce.NewClientError(err.Error(), 0, 0, cli.id)
	}

	cli.auxConnection = connection.New(auxConn, cli.settings.ResponseMessageLengthLimit, cli.id)

	return nil
}

// Stop stops the client.
func (cli *Client) Stop() (cerr *ce.CommonError) {
	cli.startStopLock.Lock()
	defer cli.startStopLock.Unlock()

	return cli.stop()
}

func (cli *Client) stop() (cerr *ce.CommonError) {
	if !cli.isWorking.Load() {
		return ce.NewClientError(ErrDoubleStopIsNotPossible, 0, 0, ClientIdNone)
	}

	cerr = cli.mainConnection.Break()
	if cerr != nil {
		return cerr
	}

	cerr = cli.auxConnection.Break()
	if cerr != nil {
		return cerr
	}

	cli.isWorking.Store(false)

	return nil
}

// Restart re-starts the client.
func (cli *Client) Restart(forcibly bool) (cerr *ce.CommonError) {
	cli.startStopLock.Lock()
	defer cli.startStopLock.Unlock()

	return cli.restart(forcibly)
}

func (cli *Client) restart(forcibly bool) (cerr *ce.CommonError) {
	if forcibly {
		_ = cli.stop()
	} else {
		cerr = cli.stop()
		if cerr != nil {
			return cerr
		}
	}

	return cli.start()
}
