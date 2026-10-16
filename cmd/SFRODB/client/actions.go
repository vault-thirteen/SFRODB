package main

import (
	"time"

	"github.com/vault-thirteen/auxie/errors"

	"github.com/vault-thirteen/SFRODB/pkg/SFRODB/classes/Client"
	ce "github.com/vault-thirteen/SFRODB/pkg/SFRODB/classes/CommonError"
	"github.com/vault-thirteen/SFRODB/pkg/common"
)

const (
	ItemLenLimit = 60
)

func makeSomeActions(cli *client.Client, appMustBeStopped *chan bool) {
	var action byte
	var uid string
	var err error
	var cerr *ce.CommonError
	var normalExit = false

	for {
		action, err = common.GetInputCharFromStdin(common.Msg_TestClientHintMain)
		if err != nil {
			errors.LogError(err)
			continue
		}

		if (action == 'q') || (action == 'Q') {
			normalExit = true
			break
		}

		switch action {
		case 'r', 'R':
		case 'g', 'G', 'e', 'E', 's', 'S', 'f', 'F':
			uid, err = common.GetInputStringFromStdin(common.Msg_TestClientHintUid)
			if err != nil {
				errors.LogError(err)
				continue
			}
		default:
			continue
		}

		switch action {
		case 'g', 'G':
			cerr = processGKeys(cli, uid)
		case 'e', 'E':
			cerr = processEKeys(cli, uid)
		case 's', 'S':
			cerr = processSKeys(cli, uid)
		case 'f', 'F':
			cerr = cli.ForgetRecord(uid)
		case 'r', 'R':
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

	common.PrintMessage(common.FormatRequestDuration(t1e.Microseconds(), itemLen))

	var mustShowData = false
	var ch byte
	var err error
	if itemLen > ItemLenLimit {
		ch, err = common.GetInputCharFromStdin(common.Msg_TestClientHintDataSize)
		if err != nil {
			return ce.NewClientError(err.Error(), 0, 0, cli.GetId())
		}

		if (ch == 'Y') || (ch == 'y') {
			mustShowData = true
		}
	} else {
		mustShowData = true
	}

	if mustShowData {
		common.PrintMessage(common.Msg_Item)
		common.PrintMessage(common.Msg_HorizontalLine)
		common.PrintMessage(string(data))
		common.PrintMessage(common.Msg_HorizontalLine)
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
		common.PrintMessage(common.Msg_RecordExists)
	} else {
		common.PrintMessage(common.Msg_RecordDoesNotExist)
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
		common.PrintMessage(common.Msg_FileExists)
	} else {
		common.PrintMessage(common.Msg_FileDoesNotExist)
	}

	return nil
}
