package settings

import (
	"errors"

	"github.com/vault-thirteen/SFRODB/common"
)

const ResponseMessageLengthLimitDefault = 1_000_000 // 1 MB.

// Settings is Client's settings.
type Settings struct {
	// Server's host name.
	Host string

	// Main port.
	MainPort uint16

	// Auxiliary port.
	AuxPort uint16

	// Maximum size for server's messages.
	ResponseMessageLengthLimit uint
}

func (stn *Settings) Check() (err error) {
	if len(stn.Host) == 0 {
		return errors.New(common.ErrClientHostIsNotSet)
	}

	if stn.MainPort == 0 {
		return errors.New(common.ErrClientPortIsNotSet)
	}

	if stn.AuxPort == 0 {
		return errors.New(common.ErrClientPortIsNotSet)
	}

	if stn.ResponseMessageLengthLimit == 0 {
		return errors.New(common.ErrResponseMessageLengthLimit)
	}

	return nil
}

func NewSettings(
	host string,
	mainPort uint16,
	auxPort uint16,
	responseMessageLengthLimit uint,
) (stn *Settings, err error) {

	stn = &Settings{
		Host:     host,
		MainPort: mainPort,
		AuxPort:  auxPort,
	}

	if responseMessageLengthLimit == 0 {
		stn.ResponseMessageLengthLimit = ResponseMessageLengthLimitDefault
	} else {
		stn.ResponseMessageLengthLimit = responseMessageLengthLimit
	}

	return stn, nil
}
