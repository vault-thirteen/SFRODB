package main

import (
	"fmt"
	"log"
	"time"

	"github.com/vault-thirteen/SFRODB/pkg/client"
	"github.com/vault-thirteen/SFRODB/pkg/common/error"
)

const (
	ErrUnsupportedKey = "unsupported key: "
	HorizontalLine    = "------------------------------------------------------------"
	ItemLenLimit      = 60
)

func makeSomeActions(cli *client.Client, appMustBeStopped *chan bool) {
	var action byte
	var tbc byte
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
		case 'g', 'G', 'e', 'E', 's', 'S', 'f', 'F', 'r', 'R':
		default:
			continue
		}

		tbc, err = getUserInputChar(HintTB)
		if err != nil {
			log.Println(err.Error())
			continue
		}

		switch tbc {
		case 't', 'T', 'b', 'B':
		case 'c', 'C':
			continue
		default:
			continue
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
			cerr = processGKeys(cli, tbc, uid)
		case 'e', 'E':
			cerr = processEKeys(cli, tbc, uid)
		case 's', 'S':
			cerr = processSKeys(cli, tbc, uid)
		case 'f', 'F':
			cerr = processFKeys(cli, tbc, uid)
		case 'r', 'R':
			cerr = processRKeys(cli, tbc)
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

func processGKeys(cli *client.Client, bt byte, uid string) (cerr *ce.CommonError) {
	var item any
	t1 := time.Now()
	item, cerr = getRecord(cli, bt, uid)
	if cerr != nil {
		return cerr
	}
	t1e := time.Now().Sub(t1)

	var itemLen int
	switch v := item.(type) {
	case string:
		itemLen = len(v)
	case []byte:
		itemLen = len(v)
	}

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
		fmt.Println(item)
		fmt.Println(HorizontalLine)
	}

	return nil
}

func getRecord(cli *client.Client, bt byte, uid string) (item any, cerr *ce.CommonError) {
	if (bt == 't') || (bt == 'T') {
		return cli.ShowText(uid)
	}

	if (bt == 'b') || (bt == 'B') {
		return cli.ShowBinary(uid)
	}

	return nil, ce.NewClientError(ErrUnsupportedKey+string(bt), 0, cli.GetId())
}

func processEKeys(cli *client.Client, bt byte, uid string) (cerr *ce.CommonError) {
	var recExists bool
	recExists, cerr = searchRecord(cli, bt, uid)
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

func searchRecord(cli *client.Client, bt byte, uid string) (exists bool, cerr *ce.CommonError) {
	if (bt == 't') || (bt == 'T') {
		return cli.SearchTextRecord(uid)
	}

	if (bt == 'b') || (bt == 'B') {
		return cli.SearchBinaryRecord(uid)
	}

	return false, ce.NewClientError(ErrUnsupportedKey+string(bt), 0, cli.GetId())
}

func processSKeys(cli *client.Client, bt byte, uid string) (cerr *ce.CommonError) {
	var fileExists bool
	fileExists, cerr = searchFile(cli, bt, uid)
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

func searchFile(cli *client.Client, bt byte, uid string) (exists bool, cerr *ce.CommonError) {
	if (bt == 't') || (bt == 'T') {
		return cli.SearchTextFile(uid)
	}

	if (bt == 'b') || (bt == 'B') {
		return cli.SearchBinaryFile(uid)
	}

	return false, ce.NewClientError(ErrUnsupportedKey+string(bt), 0, cli.GetId())
}

func processFKeys(cli *client.Client, bt byte, uid string) (cerr *ce.CommonError) {
	return forgetRecord(cli, bt, uid)
}

func forgetRecord(cli *client.Client, bt byte, uid string) (cerr *ce.CommonError) {
	if (bt == 't') || (bt == 'T') {
		return cli.ForgetTextRecord(uid)
	}

	if (bt == 'b') || (bt == 'B') {
		return cli.ForgetBinaryRecord(uid)
	}

	return ce.NewClientError(ErrUnsupportedKey+string(bt), 0, cli.GetId())
}

func processRKeys(cli *client.Client, bt byte) (cerr *ce.CommonError) {
	return resetCache(cli, bt)
}

func resetCache(cli *client.Client, bt byte) (cerr *ce.CommonError) {
	if (bt == 't') || (bt == 'T') {
		return cli.ResetTextCache()
	}

	if (bt == 'b') || (bt == 'B') {
		return cli.ResetBinaryCache()
	}

	return ce.NewClientError(ErrUnsupportedKey+string(bt), 0, cli.GetId())
}
