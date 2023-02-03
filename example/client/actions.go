package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/vault-thirteen/SFRODB/client"
)

const (
	ErrUnknownAction = "unknown action: "
	HintGet          = "To get a textual item – enter T, B – for binary item, Q – to quit.\r\n>"
	HintUid          = "Enter the UID.\r\n>"
	HintDataSize     = `"Data is quite large. Do you want to see it ?
Enter 'Y' for Yes, otherwise – No.
>`
	HorizontalLine = "------------------------------------------------------------"
	ItemLenLimit   = 60
)

func makeSomeActions(cli *client.Client, appMustBeStopped *chan bool) {
	var ch byte
	var err error
	var s string
	var item any
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

		if (ch != 't') && (ch != 'T') &&
			(ch != 'b') && (ch != 'B') {
			continue
		}

		s, err = getUserInputString(HintUid)
		if err != nil {
			log.Println(err.Error())
			continue
		}

		t1 := time.Now()
		item, err = getItem(cli, ch, s)
		if err != nil {
			log.Println(err.Error())
			break
		}
		t1e := time.Now().Sub(t1)

		var itemLen int
		switch v := item.(type) {
		case string:
			itemLen = len(v)
		case []byte:
			itemLen = len(v)
		}

		fmt.Printf("Request Duration: %d µs. Data Size: %d Bytes.\r\n", t1e.Microseconds(), itemLen)

		var mustShowData = false
		if itemLen > ItemLenLimit {
			ch, err = getUserInputChar(HintDataSize)
			if err != nil {
				log.Println(err.Error())
				continue
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
	}

	// Send thr 'Close' Request.
	err = cli.SayGoodbye(normalExit)
	if err != nil {
		log.Println(err.Error())
	}

	*appMustBeStopped <- true
}

func getUserInputChar(hint string) (char byte, err error) {
	if len(hint) > 0 {
		fmt.Print(hint)
	}

	reader := bufio.NewReader(os.Stdin)

	char, err = reader.ReadByte()
	if err != nil {
		return 0, err
	}

	return char, nil
}

func getUserInputString(hint string) (s string, err error) {
	if len(hint) > 0 {
		fmt.Print(hint)
	}

	reader := bufio.NewReader(os.Stdin)

	s, err = reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(s), nil
}

func getItem(cli *client.Client, action byte, uid string) (item any, err error) {
	if (action == 't') || (action == 'T') {
		var text string
		text, err = cli.GetText(uid)
		if err != nil {
			return nil, err
		}

		return text, nil
	}

	if (action == 'b') || (action == 'B') {
		var data []byte
		data, err = cli.GetBinary(uid)
		if err != nil {
			return nil, err
		}

		return data, nil
	}

	return nil, errors.New(ErrUnknownAction + string(action))
}
