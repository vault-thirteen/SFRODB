package settings

import (
	"errors"

	"github.com/vault-thirteen/SFRODB/common"
)

const ResponseMessageLengthLimitDefault = 1000000 // 1 MB.

type Settings struct {
	// Client's host name.
	ClientHost string

	// Client's port.
	ClientPort uint16

	// Maximum size for server's messages.
	ResponseMessageLengthLimit uint
}

func (stn *Settings) Check() (err error) {
	if len(stn.ClientHost) == 0 {
		return errors.New(common.ErrClientHostIsNotSet)
	}

	if stn.ClientPort == 0 {
		return errors.New(common.ErrClientPortIsNotSet)
	}

	if stn.ResponseMessageLengthLimit == 0 {
		return errors.New(common.ErrResponseMessageLengthLimit)
	}

	return nil
}

func NewSettings(
	host string,
	port uint16,
	responseMessageLengthLimit uint,
) (stn *Settings, err error) {

	stn = &Settings{
		ClientHost: host,
		ClientPort: port,
	}

	if responseMessageLengthLimit == 0 {
		stn.ResponseMessageLengthLimit = ResponseMessageLengthLimitDefault
	} else {
		stn.ResponseMessageLengthLimit = responseMessageLengthLimit
	}

	return stn, nil
}
