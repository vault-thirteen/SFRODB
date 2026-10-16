package common

import (
	"bufio"
	"os"
	"strings"
)

func GetInputCharFromStdin(hint string) (char byte, err error) {
	if len(hint) > 0 {
		PrintMessage(hint)
	}

	reader := bufio.NewReader(os.Stdin)
	return reader.ReadByte()
}

func GetInputStringFromStdin(hint string) (s string, err error) {
	if len(hint) > 0 {
		PrintMessage(hint)
	}

	reader := bufio.NewReader(os.Stdin)

	s, err = reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(s), nil
}
