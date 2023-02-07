package main

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/vault-thirteen/SFRODB/client"
	ce "github.com/vault-thirteen/SFRODB/common/error"
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
			err = processGKeys(cli, tbc, uid)
		case 'e', 'E':
			err = processEKeys(cli, tbc, uid)
		case 's', 'S':
			err = processSKeys(cli, tbc, uid)
		case 'f', 'F':
			err = processFKeys(cli, tbc, uid)
		case 'r', 'R':
			err = processRKeys(cli, tbc)
		default:
			continue
		}
		if err != nil {
			de := getDetailedError(err)
			if de != nil {
				if de.IsServerError() {
					log.Println("Server Error: " + err.Error())
					break
				} else if de.IsClientError() {
					log.Println("Client Error: " + err.Error())
					continue
				} else {
					log.Println("Anomaly: " + err.Error())
					break
				}
			} else {
				log.Println(err.Error())
				continue
			}
		}
	}

	// Send the 'Close' Request.
	err = cli.CloseConnection_Main(normalExit)
	if err != nil {
		log.Println(err.Error())
	}

	err = cli.CloseConnection_Aux(normalExit)
	if err != nil {
		log.Println(err.Error())
	}

	*appMustBeStopped <- true
}

func getDetailedError(err error) (de *ce.CommonError) {
	detailedError, ok := err.(*ce.CommonError)
	if !ok {
		return nil
	}

	return detailedError
}

func processGKeys(cli *client.Client, bt byte, uid string) (err error) {
	var item any
	t1 := time.Now()
	item, err = getRecord(cli, bt, uid)
	if err != nil {
		return err
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
	if itemLen > ItemLenLimit {
		ch, err = getUserInputChar(HintDataSize)
		if err != nil {
			return err
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

func getRecord(cli *client.Client, bt byte, uid string) (item any, err error) {
	if (bt == 't') || (bt == 'T') {
		return cli.ShowText(uid)
	}

	if (bt == 'b') || (bt == 'B') {
		return cli.ShowBinary(uid)
	}

	return nil, errors.New(ErrUnsupportedKey + string(bt))
}

func processEKeys(cli *client.Client, bt byte, uid string) (err error) {
	var recExists bool
	recExists, err = searchRecord(cli, bt, uid)
	if err != nil {
		return err
	}

	if recExists {
		fmt.Println("Record exists.")
	} else {
		fmt.Println("Record does not exist.")
	}

	return nil
}

func searchRecord(cli *client.Client, bt byte, uid string) (exists bool, err error) {
	if (bt == 't') || (bt == 'T') {
		return cli.SearchTextRecord(uid)
	}

	if (bt == 'b') || (bt == 'B') {
		return cli.SearchBinaryRecord(uid)
	}

	return false, errors.New(ErrUnsupportedKey + string(bt))
}

func processSKeys(cli *client.Client, bt byte, uid string) (err error) {
	var fileExists bool
	fileExists, err = searchFile(cli, bt, uid)
	if err != nil {
		return err
	}

	if fileExists {
		fmt.Println("File exists.")
	} else {
		fmt.Println("File does not exist.")
	}

	return nil
}

func searchFile(cli *client.Client, bt byte, uid string) (exists bool, err error) {
	if (bt == 't') || (bt == 'T') {
		return cli.SearchTextFile(uid)
	}

	if (bt == 'b') || (bt == 'B') {
		return cli.SearchBinaryFile(uid)
	}

	return false, errors.New(ErrUnsupportedKey + string(bt))
}

func processFKeys(cli *client.Client, bt byte, uid string) (err error) {
	return forgetRecord(cli, bt, uid)
}

func forgetRecord(cli *client.Client, bt byte, uid string) (err error) {
	if (bt == 't') || (bt == 'T') {
		return cli.ForgetTextRecord(uid)
	}

	if (bt == 'b') || (bt == 'B') {
		return cli.ForgetBinaryRecord(uid)
	}

	return errors.New(ErrUnsupportedKey + string(bt))
}

func processRKeys(cli *client.Client, bt byte) (err error) {
	return resetCache(cli, bt)
}

func resetCache(cli *client.Client, bt byte) (err error) {
	if (bt == 't') || (bt == 'T') {
		return cli.ResetTextCache()
	}

	if (bt == 'b') || (bt == 'B') {
		return cli.ResetBinaryCache()
	}

	return errors.New(ErrUnsupportedKey + string(bt))
}
