package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	CRLF     = "\r\n"
	HintMain = CRLF +
		"[G] = Get/Show a Record;" + CRLF +
		"[E] = Check Record's Existence;" + CRLF +
		"[S] = Check File's Existence;" + CRLF +
		"[F] = Forget/Remove a Record from Cache;" + CRLF +
		"[R] = Reset/Clear the Cache;" + CRLF +
		"[Q] = Quit/Exit." + CRLF +
		"> "
	HintUid      = "Enter the UID > "
	HintDataSize = "Data is quite large. Do you want to see it ? [Y] = Yes; [N] = No. > "
)

func getUserInputChar(hint string) (char byte, err error) {
	if len(hint) > 0 {
		fmt.Print(hint)
	}

	reader := bufio.NewReader(os.Stdin)
	return reader.ReadByte()
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
