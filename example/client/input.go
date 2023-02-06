package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	HintMain = "\r\n" +
		"[G] = Get/Show a Record;\r\n" +
		"[E] = Check Record's Existence;\r\n" +
		"[S] = Check File's Existence;\r\n" +
		"[F] = Forget/Remove a Record from Cache;\r\n" +
		"[R] = Reset/Clear the Cache;\r\n" +
		"[Q] = Quit/Exit.\r\n> "
	HintTB       = "[T] = Text; [B] = Binary; [C] = Cancel/Quit. \r\n> "
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
