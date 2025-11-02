package cs

import (
	"errors"
)

const ResponseMessageLengthLimitDefault = 1_000_000 // 1 MB.

const (
	ErrClientHostIsNotSet         = "client host is not set"
	ErrClientPortIsNotSet         = "client port is not set"
	ErrResponseMessageLengthLimit = "response message length limit is not set"
)

// ClientSettings is client's Settings.
type ClientSettings struct {
	// Server's host name.
	Host string

	// Main port.
	MainPort uint16

	// Auxiliary port.
	AuxPort uint16

	// Maximum size for server's messages.
	ResponseMessageLengthLimit uint
}

func New(
	host string,
	mainPort uint16,
	auxPort uint16,
	responseMessageLengthLimit uint,
) (stn *ClientSettings, err error) {
	stn = &ClientSettings{
		Host:     host,
		MainPort: mainPort,
		AuxPort:  auxPort,
	}

	stn.ResponseMessageLengthLimit = responseMessageLengthLimit

	if stn.ResponseMessageLengthLimit == 0 {
		stn.ResponseMessageLengthLimit = ResponseMessageLengthLimitDefault
	}

	return stn, nil
}

func (stn *ClientSettings) Check() (err error) {
	if len(stn.Host) == 0 {
		return errors.New(ErrClientHostIsNotSet)
	}

	if stn.MainPort == 0 {
		return errors.New(ErrClientPortIsNotSet)
	}

	if stn.AuxPort == 0 {
		return errors.New(ErrClientPortIsNotSet)
	}

	if stn.ResponseMessageLengthLimit == 0 {
		return errors.New(ErrResponseMessageLengthLimit)
	}

	return nil
}
