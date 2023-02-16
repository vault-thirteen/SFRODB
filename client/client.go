package client

import (
	"fmt"
	"net"
	"sync"
	"sync/atomic"

	"github.com/vault-thirteen/SFRODB/client/settings"
	"github.com/vault-thirteen/SFRODB/common"
	"github.com/vault-thirteen/SFRODB/common/connection"
	ce "github.com/vault-thirteen/SFRODB/common/error"
	"github.com/vault-thirteen/SFRODB/common/method"
	"github.com/vault-thirteen/SFRODB/common/protocol"
	"github.com/vault-thirteen/errorz"
)

const (
	// ClientIdNone is a client ID for a non-network operations.
	ClientIdNone = ""

	// ClientIdIncoming is a client ID for a request incoming to server.
	ClientIdIncoming = "s"
)

const (
	ErrDoubleStartIsNotPossible = "double start is not possible"
	ErrDoubleStopIsNotPossible  = "double stop is not possible"
)

// Client is client.
type Client struct {
	// All the clients must have unique IDs.
	id string

	settings *settings.Settings

	mainDsn string
	auxDsn  string

	mainAddr *net.TCPAddr
	auxAddr  *net.TCPAddr

	methodNameBuffers map[method.Method][]byte
	methodValues      map[string]method.Method

	mainConnection *connection.Connection
	auxConnection  *connection.Connection

	// Internal control structures.
	startStopLock *sync.Mutex
	isWorking     *atomic.Bool
}

// NewClient creates a client.
func NewClient(stn *settings.Settings, id string) (cli *Client, err error) {
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
		return ce.NewClientError(ErrDoubleStartIsNotPossible, 0, ClientIdNone)
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
	var mainConn *net.TCPConn
	var err error
	mainConn, err = net.DialTCP(proto.LowLevelProtocol, nil, cli.mainAddr)
	if err != nil {
		return ce.NewClientError(err.Error(), 0, cli.id)
	}

	err = connection.EnableKeepAlives(mainConn)
	if err != nil {
		closeErr := mainConn.Close()
		if closeErr != nil {
			err = errorz.Combine(err, closeErr)
		}
		return ce.NewClientError(err.Error(), 0, cli.id)
	}

	cli.mainConnection = connection.NewConnection(
		mainConn,
		&cli.methodNameBuffers,
		&cli.methodValues,
		cli.settings.ResponseMessageLengthLimit,
		cli.id,
	)

	return nil
}

func (cli *Client) startAuxConnection() (cerr *ce.CommonError) {
	var auxConn *net.TCPConn
	var err error
	auxConn, err = net.DialTCP(proto.LowLevelProtocol, nil, cli.auxAddr)
	if err != nil {
		return ce.NewClientError(err.Error(), 0, cli.id)
	}

	err = connection.EnableKeepAlives(auxConn)
	if err != nil {
		closeErr := auxConn.Close()
		if closeErr != nil {
			err = errorz.Combine(err, closeErr)
		}
		return ce.NewClientError(err.Error(), 0, cli.id)
	}

	cli.auxConnection = connection.NewConnection(
		auxConn,
		&cli.methodNameBuffers,
		&cli.methodValues,
		cli.settings.ResponseMessageLengthLimit,
		cli.id,
	)

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
		return ce.NewClientError(ErrDoubleStopIsNotPossible, 0, ClientIdNone)
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

// LibVersion returns the version of this library.
func (cli *Client) LibVersion() (ver string) {
	return common.LibVersion
}
