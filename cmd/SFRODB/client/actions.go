package main

import (
	"fmt"
	"log"
	"time"

	"github.com/vault-thirteen/SFRODB/pkg/SFRODB/classes/Client"
	ce "github.com/vault-thirteen/SFRODB/pkg/SFRODB/classes/CommonError"
	"github.com/vault-thirteen/SFRODB/pkg/common"
)

const (
	HorizontalLine = "------------------------------------------------------------"
	ItemLenLimit   = 60
)

const (
	Msgf_RequestDuration   = "Request Duration: %d µs. Data Size: %d Bytes."
	Msg_Item               = "Item:"
	Msg_RecordExists       = "Record exists."
	Msg_RecordDoesNotExist = "Record does not exist."
	Msg_FileExists         = "File exists."
	Msg_FileDoesNotExist   = "File does not exist."
)

func makeSomeActions(cli *client.Client, appMustBeStopped *chan bool) {
	var action byte
	var uid string
	var err error
	var cerr *ce.CommonError
	var normalExit = false

	for {
		action, err = getUserInputChar(HintMain)
		if err != nil {
			log.Println(err.Error())
			continue
		}

		if (action == 'q') || (action == 'Q') {
			normalExit = true
			break
		}

		switch action {
		case 'r', 'R':
		case 'g', 'G', 'e', 'E', 's', 'S', 'f', 'F':
			uid, err = getUserInputString(HintUid)
			if err != nil {
				log.Println(err.Error())
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
				log.Println(common.Msg_ServerError + cerr.Error())
				break
			} else if cerr.IsClientError() {
				log.Println(common.Msg_ClientError + cerr.Error())
				continue
			} else {
				log.Println(common.Msg_Anomaly + cerr.Error())
				break
			}
		}
	}

	// Send the 'Close' Request.
	cerr = cli.CloseConnection_Main(normalExit)
	if cerr != nil {
		log.Println(cerr.Error())
	}

	cerr = cli.CloseConnection_Aux(normalExit)
	if cerr != nil {
		log.Println(cerr.Error())
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

	fmt.Println(fmt.Sprintf(Msgf_RequestDuration, t1e.Microseconds(), itemLen))

	var mustShowData = false
	var ch byte
	var err error
	if itemLen > ItemLenLimit {
		ch, err = getUserInputChar(HintDataSize)
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
		fmt.Println(Msg_Item)
		fmt.Println(HorizontalLine)
		fmt.Println(string(data))
		fmt.Println(HorizontalLine)
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
		fmt.Println(Msg_RecordExists)
	} else {
		fmt.Println(Msg_RecordDoesNotExist)
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
		fmt.Println(Msg_FileExists)
	} else {
		fmt.Println(Msg_FileDoesNotExist)
	}

	return nil
}
