package main

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/vault-thirteen/SFRODB/client"
)

const (
	ErrUnknownAction = "unknown action: "
	HorizontalLine   = "------------------------------------------------------------"
	ItemLenLimit     = 60
)

func makeSomeActions(cli *client.Client, appMustBeStopped *chan bool) {
	var ch byte
	var err error
	var normalExit = false

	for {
		ch, err = getUserInputChar(HintGet)
		if err != nil {
			log.Println(err.Error())
			continue
		}

		if (ch == 'q') || (ch == 'Q') {
			normalExit = true
			break
		}

		if (ch == 'b') || (ch == 'B') ||
			(ch == 't') || (ch == 'T') {
			err = processBTKeys(cli, ch)
			if err != nil {
				log.Println(err.Error())
				break
			}
		}

		if (ch == 'x') || (ch == 'X') {
			err = processXKeys(cli)
			if err != nil {
				log.Println(err.Error())
				break
			}
		}
	}

	// Send the 'Close' Request.
	err = cli.SayGoodbyeOnMain(normalExit)
	if err != nil {
		log.Println(err.Error())
	}

	err = cli.SayGoodbyeOnAux(normalExit)
	if err != nil {
		log.Println(err.Error())
	}

	*appMustBeStopped <- true
}

func processBTKeys(cli *client.Client, ch byte) (err error) {
	var s string
	s, err = getUserInputString(HintUid)
	if err != nil {
		return err
	}

	var item any
	t1 := time.Now()
	item, err = getItem(cli, ch, s)
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

func processXKeys(cli *client.Client) (err error) {
	var ch byte
	ch, err = getUserInputChar(HintExtra)
	if err != nil {
		return err
	}

	if (ch == 'b') || (ch == 'B') ||
		(ch == 't') || (ch == 'T') {
		return processXBTKeys(cli, ch)
	}

	if (ch == 'c') || (ch == 'C') {
		return processXCKeys(cli)
	}

	return nil
}

func processXBTKeys(cli *client.Client, ch byte) (err error) {
	var s string
	s, err = getUserInputString(HintUid)
	if err != nil {
		return err
	}

	return removeItem(cli, ch, s)
}

func processXCKeys(cli *client.Client) (err error) {
	var ch byte
	ch, err = getUserInputChar(HintExtraC)
	if err != nil {
		return err
	}

	if (ch == 'b') || (ch == 'B') ||
		(ch == 't') || (ch == 'T') {
		return clearCache(cli, ch)
	}

	return nil
}

func getItem(cli *client.Client, action byte, uid string) (item any, err error) {
	if (action == 't') || (action == 'T') {
		return cli.GetText(uid)
	}

	if (action == 'b') || (action == 'B') {
		return cli.GetBinary(uid)
	}

	return nil, errors.New(ErrUnknownAction + string(action))
}

func removeItem(cli *client.Client, action byte, uid string) (err error) {
	if (action == 't') || (action == 'T') {
		return cli.RemoveText(uid)
	}

	if (action == 'b') || (action == 'B') {
		return cli.RemoveBinary(uid)
	}

	return errors.New(ErrUnknownAction + string(action))
}

func clearCache(cli *client.Client, action byte) (err error) {
	if (action == 't') || (action == 'T') {
		return cli.ClearTextCache()
	}

	if (action == 'b') || (action == 'B') {
		return cli.ClearBinaryCache()
	}

	return errors.New(ErrUnknownAction + string(action))
}
