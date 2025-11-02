package ss

import (
	"errors"
	"os"
	"strings"

	ds "github.com/vault-thirteen/SFRODB/pkg/SFRODB/classes/DataSettings"
	ae "github.com/vault-thirteen/auxie/errors"
	"github.com/vault-thirteen/auxie/number"
	"github.com/vault-thirteen/auxie/reader"
)

const (
	ErrFileIsNotSet       = "file is not set"
	ErrServerHostIsNotSet = "server host is not set"
	ErrServerPortIsNotSet = "server port is not set"
)

// ServerSettings is Server's Settings.
type ServerSettings struct {
	// Path to the Settings file.
	// ServerSettings are positional for simplicity.
	File string

	// Server host name.
	Hostname string

	// Main port.
	// A port which is used for read-only operations.
	MainPort uint16

	// Auxiliary port.
	// A port which is used for auxiliary operations.
	AuxPort uint16

	// Data Settings.
	Data *ds.DataSettings
}

func NewSettingsFromFile(filePath string) (stn *ServerSettings, err error) {
	stn = &ServerSettings{
		File: filePath,
	}

	var file *os.File
	file, err = os.Open(filePath)
	if err != nil {
		return stn, err
	}
	defer func() {
		derr := file.Close()
		if derr != nil {
			err = ae.Combine(err, derr)
		}
	}()

	rdr := reader.New(file)
	var buf = make([][]byte, 5)

	for i := range buf {
		buf[i], err = rdr.ReadLineEndingWithCRLF()
		if err != nil {
			return stn, err
		}
	}

	// Server Hostname & Port.
	stn.Hostname = strings.TrimSpace(string(buf[0]))

	stn.MainPort, err = number.ParseUint16(strings.TrimSpace(string(buf[1])))
	if err != nil {
		return stn, err
	}

	stn.AuxPort, err = number.ParseUint16(strings.TrimSpace(string(buf[2])))
	if err != nil {
		return stn, err
	}

	// Cache Data ServerSettings.
	stn.Data, err = ds.ParseDataSettings(string(buf[3]), string(buf[4]))
	if err != nil {
		return stn, err
	}

	return stn, nil
}

func (stn *ServerSettings) Check() (err error) {
	if len(stn.File) == 0 {
		return errors.New(ErrFileIsNotSet)
	}

	if len(stn.Hostname) == 0 {
		return errors.New(ErrServerHostIsNotSet)
	}

	if stn.MainPort == 0 {
		return errors.New(ErrServerPortIsNotSet)
	}

	if stn.AuxPort == 0 {
		return errors.New(ErrServerPortIsNotSet)
	}

	err = stn.Data.Check()
	if err != nil {
		return err
	}

	return nil
}
