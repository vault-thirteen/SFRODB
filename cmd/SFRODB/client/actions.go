package main

import (
	"fmt"
	"log"
	"time"

	"github.com/vault-thirteen/SFRODB/pkg/SFRODB/client"
	ce "github.com/vault-thirteen/SFRODB/pkg/SFRODB/common/error"
)

const (
	HorizontalLine = "------------------------------------------------------------"
	ItemLenLimit   = 60
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
				log.Println("Server Error: " + cerr.Error())
				break
			} else if cerr.IsClientError() {
				log.Println("Client Error: " + cerr.Error())
				continue
			} else {
				log.Println("Anomaly: " + cerr.Error())
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

	fmt.Printf("Request Duration: %d µs. Data Size: %d Bytes.\r\n",
		t1e.Microseconds(), itemLen)

	var mustShowData = false
	var ch byte
	var err error
	if itemLen > ItemLenLimit {
		ch, err = getUserInputChar(HintDataSize)
		if err != nil {
			return ce.NewClientError(err.Error(), 0, cli.GetId())
		}

		if (ch == 'Y') || (ch == 'y') {
			mustShowData = true
		}
	} else {
		mustShowData = true
	}

	if mustShowData {
		fmt.Println("Item:")
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
		fmt.Println("Record exists.")
	} else {
		fmt.Println("Record does not exist.")
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
		fmt.Println("File exists.")
	} else {
		fmt.Println("File does not exist.")
	}

	return nil
}
