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

func (cli *Client) GetMainDsn() (dsn string) {
	return cli.mainDsn
}

func (cli *Client) GetAuxDsn() (dsn string) {
	return cli.auxDsn
}

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

func (cli *Client) GetText(uid string) (text string, err error) {
	var rm *common.Request
	rm, err = common.NewRequest_ShowText(uid)
	if err != nil {
		return "", err
	}

	err = cli.mainConnection.SendRequestMessage(rm)
	if err != nil {
		return "", err
	}

	var resp *common.Response
	resp, err = cli.mainConnection.GetResponseMessage()
	if err != nil {
		return "", err
	}

	// If something goes wrong, server warns about closing the mainConnection.
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

	err = cli.mainConnection.SendRequestMessage(rm)
	if err != nil {
		return nil, err
	}

	var resp *common.Response
	resp, err = cli.mainConnection.GetResponseMessage()
	if err != nil {
		return nil, err
	}

	// If something goes wrong, server warns about closing the mainConnection.
	if resp.Method == common.MethodClosingConnection {
		return nil, errors.New(common.ErrSomethingWentWrong)
	}

	return resp.Data, nil
}

func (cli *Client) SayGoodbyeOnMain(normalExit bool) (err error) {
	return cli.sayGoodbye(true, normalExit)
}

func (cli *Client) SayGoodbyeOnAux(normalExit bool) (err error) {
	return cli.sayGoodbye(false, normalExit)
}

func (cli *Client) sayGoodbye(useMainConnection bool, normalExit bool) (err error) {
	var rm *common.Request
	rm, err = common.NewRequest_CloseConnection()
	if err != nil {
		return err
	}

	if useMainConnection {
		err = cli.mainConnection.SendRequestMessage(rm)
	} else {
		err = cli.auxConnection.SendRequestMessage(rm)
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
		resp, err = cli.mainConnection.GetResponseMessage()
	} else {
		resp, err = cli.auxConnection.GetResponseMessage()
	}
	if err != nil {
		return err
	}

	if resp.Method != common.MethodClosingConnection {
		return errors.New(common.ErrSomethingWentWrong)
	}

	return nil
}

func (cli *Client) RemoveText(uid string) (err error) {
	var rm *common.Request
	rm, err = common.NewRequest_RemoveText(uid)
	if err != nil {
		return err
	}

	err = cli.auxConnection.SendRequestMessage(rm)
	if err != nil {
		return err
	}

	var resp *common.Response
	resp, err = cli.auxConnection.GetResponseMessage()
	if err != nil {
		return err
	}

	if resp.Method != common.MethodOK {
		return errors.New(common.ErrSomethingWentWrong)
	}

	return nil
}

func (cli *Client) RemoveBinary(uid string) (err error) {
	var rm *common.Request
	rm, err = common.NewRequest_RemoveBinary(uid)
	if err != nil {
		return err
	}

	err = cli.auxConnection.SendRequestMessage(rm)
	if err != nil {
		return err
	}

	var resp *common.Response
	resp, err = cli.auxConnection.GetResponseMessage()
	if err != nil {
		return err
	}

	if resp.Method != common.MethodOK {
		return errors.New(common.ErrSomethingWentWrong)
	}

	return nil
}

func (cli *Client) ClearTextCache() (err error) {
	var rm *common.Request
	rm, err = common.NewRequest_ClearTextCache()
	if err != nil {
		return err
	}

	err = cli.auxConnection.SendRequestMessage(rm)
	if err != nil {
		return err
	}

	var resp *common.Response
	resp, err = cli.auxConnection.GetResponseMessage()
	if err != nil {
		return err
	}

	if resp.Method != common.MethodOK {
		return errors.New(common.ErrSomethingWentWrong)
	}

	return nil
}

func (cli *Client) ClearBinaryCache() (err error) {
	var rm *common.Request
	rm, err = common.NewRequest_ClearBinaryCache()
	if err != nil {
		return err
	}

	err = cli.auxConnection.SendRequestMessage(rm)
	if err != nil {
		return err
	}

	var resp *common.Response
	resp, err = cli.auxConnection.GetResponseMessage()
	if err != nil {
		return err
	}

	if resp.Method != common.MethodOK {
		return errors.New(common.ErrSomethingWentWrong)
	}

	return nil
}
