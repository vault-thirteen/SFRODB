package main

import (
	"errors"
	"os"

	"github.com/vault-thirteen/SFRODB/pkg/common/helper"
)

const ErrSyntax = "syntax error"

type CommandLineArguments struct {
	Host     string
	MainPort uint16
	AuxPort  uint16
}

func readCLA() (cla *CommandLineArguments, err error) {
	if len(os.Args) != 4 {
		return nil, errors.New(ErrSyntax)
	}

	cla = &CommandLineArguments{
		Host: os.Args[1],
	}

	cla.MainPort, err = helper.ParseUint16(os.Args[2])
	if err != nil {
		return nil, err
	}

	cla.AuxPort, err = helper.ParseUint16(os.Args[3])
	if err != nil {
		return nil, err
	}

	return cla, nil
}
