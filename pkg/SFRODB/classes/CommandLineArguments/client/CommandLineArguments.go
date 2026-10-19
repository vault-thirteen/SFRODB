package clac

import (
	"os"

	ae "github.com/vault-thirteen/auxie/errors"
	"github.com/vault-thirteen/auxie/number"

	"github.com/vault-thirteen/SFRODB/pkg/common"
)

type CommandLineArguments struct {
	Host     string
	MainPort uint16
	AuxPort  uint16
}

func NewCommandLineArgumentsFromOs() (cla *CommandLineArguments, err error) {
	if len(os.Args) != 4 {
		return nil, ae.ErrorFromText(common.Err_Syntax)
	}

	cla = &CommandLineArguments{
		Host: os.Args[1],
	}

	cla.MainPort, err = number.ParseUint16(os.Args[2])
	if err != nil {
		return nil, err
	}

	cla.AuxPort, err = number.ParseUint16(os.Args[3])
	if err != nil {
		return nil, err
	}

	return cla, nil
}
