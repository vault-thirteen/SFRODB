package main

import (
	"errors"
	"os"

	"github.com/vault-thirteen/SFRODB/common"
)

const ErrSyntax = "syntax error"

type CommandLineArguments struct {
	Host string
	Port uint16
}

func readCLA() (cla *CommandLineArguments, err error) {
	if len(os.Args) != 3 {
		return nil, errors.New(ErrSyntax)
	}

	cla = &CommandLineArguments{
		Host: os.Args[1],
	}

	cla.Port, err = common.ParseUint16(os.Args[2])
	if err != nil {
		return nil, err
	}

	return cla, nil
}
