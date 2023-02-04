package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	HintGet      = "To get a textual item – enter T, B – for binary item, X - extras, Q – to quit.\r\n>"
	HintUid      = "Enter the UID.\r\n>"
	HintDataSize = `"Data is quite large. Do you want to see it ?
Enter 'Y' for Yes, otherwise – No.
>`
	HintExtra  = "To remove a textual item – enter T, B – for binary item, C - to clear all the cache.\r\n>"
	HintExtraC = "To clear the textual cache – enter T, B – for binary cache.\r\n>"
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
