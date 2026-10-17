package main

import (
	"time"

	"github.com/vault-thirteen/auxie/cio"
	"github.com/vault-thirteen/auxie/errors"

	"github.com/vault-thirteen/SFRODB/pkg/SFRODB/classes/Client"
	ce "github.com/vault-thirteen/SFRODB/pkg/SFRODB/classes/CommonError"
	"github.com/vault-thirteen/SFRODB/pkg/common"
)

const (
	ItemLenLimit = 60
)

func makeSomeActions(cli *client.Client, appMustBeStopped *chan bool) {
	var err error
	err = cio.EnableAnsiSymbolsInTerminal()
	if err != nil {
		return
	}

	var key cio.Key
	var uid string
	var cerr *ce.CommonError
	var normalExit = false

	for {
		cio.Print(common.Msg_TestClientHintMain)
		cio.Print(common.Msg_EraseCurrentLine)
		key, err = cio.GetInputKey()
		if err != nil {
			errors.LogError(err)
			continue
		}

		if key.Code() == cio.VK_Q {
			normalExit = true
			break
		}

		switch key.Code() {
		case cio.VK_R:
		case cio.VK_G,
			cio.VK_E,
			cio.VK_S,
			cio.VK_F:
			{
				cio.Print(common.Msg_TestClientHintUid)

				uid, err = cio.GetInputString(true)
				if err != nil {
					errors.LogError(err)
					continue
				}
			}
		default:
			continue
		}

		switch key.Code() {
		case cio.VK_G:
			cerr = processGKeys(cli, uid)
		case cio.VK_E:
			cerr = processEKeys(cli, uid)
		case cio.VK_S:
			cerr = processSKeys(cli, uid)
		case cio.VK_F:
			cerr = cli.ForgetRecord(uid)
		case cio.VK_R:
			cerr = cli.ResetCache()
		default:
			continue
		}
		if cerr != nil {
			if cerr.IsServerError() {
				errors.LogError(errors.ErrorWithPrefix(cerr, common.Err_ServerError))
				break
			} else if cerr.IsClientError() {
				errors.LogError(errors.ErrorWithPrefix(cerr, common.Err_ClientError))
				continue
			} else {
				errors.LogError(errors.ErrorWithPrefix(cerr, common.Err_Anomaly))
				break
			}
		}
	}

	// Send the 'Close' Request.
	cerr = cli.CloseConnection_Main(normalExit)
	if cerr != nil {
		errors.LogError(cerr)
	}

	cerr = cli.CloseConnection_Aux(normalExit)
	if cerr != nil {
		errors.LogError(cerr)
	}

	*appMustBeStopped <- true
}

func processGKeys(cli *client.Client, uid string) (cerr *ce.CommonError) {
	t1 := time.Now()

	var data []byte
	data, cerr = cli.ShowData(uid)
	if cerr != nil {
		return cerr
	}
	var itemLen = len(data)

	t1e := time.Now().Sub(t1)

	cio.PrintMessage(common.FormatRequestDuration(t1e.Microseconds(), itemLen))

	var mustShowData = false
	var key cio.Key
	var err error
	if itemLen > ItemLenLimit {
		cio.Print(common.Msg_TestClientHintDataSize)
		key, err = cio.GetInputKey()
		if err != nil {
			return ce.NewClientError(err.Error(), 0, 0, cli.GetId())
		}

		if key.Code() == cio.VK_Y {
			mustShowData = true
		}
	} else {
		mustShowData = true
	}

	if mustShowData {
		cio.PrintMessage(common.Msg_Item)
		cio.PrintMessage(common.Msg_HorizontalLine)
		cio.PrintMessage(string(data))
		cio.PrintMessage(common.Msg_HorizontalLine)
	}

	return nil
}

func processEKeys(cli *client.Client, uid string) (cerr *ce.CommonError) {
	var recExists bool
	recExists, cerr = cli.SearchRecord(uid)
	if cerr != nil {
		return cerr
	}

	if recExists {
		cio.PrintMessage(common.Msg_RecordExists)
	} else {
		cio.PrintMessage(common.Msg_RecordDoesNotExist)
	}

	return nil
}

func processSKeys(cli *client.Client, uid string) (cerr *ce.CommonError) {
	var fileExists bool
	fileExists, cerr = cli.SearchFile(uid)
	if cerr != nil {
		return cerr
	}

	if fileExists {
		cio.PrintMessage(common.Msg_FileExists)
	} else {
		cio.PrintMessage(common.Msg_FileDoesNotExist)
	}

	return nil
}
