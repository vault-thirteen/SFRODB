package settings

import (
	"errors"
	"os"
	"strings"

	"github.com/vault-thirteen/SFRODB/pkg/common/error"
	"github.com/vault-thirteen/auxie/number"
	"github.com/vault-thirteen/auxie/reader"
	"github.com/vault-thirteen/errorz"
)

// Settings is Server's settings.
type Settings struct {
	// Path to the settings file.
	// Settings are positional for simplicity.
	File string

	// Host name.
	ServerHost string

	// Main port.
	// A port which is used for read-only operations.
	MainPort uint16

	// Auxiliary port.
	// A port which is used for non-read operations.
	AuxPort uint16

	// Data settings.
	Data *DataSettings
}

func NewSettingsFromFile(filePath string) (stn *Settings, err error) {
	stn = &Settings{
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
			err = errorz.Combine(err, derr)
		}
	}()

	rdr := reader.NewReader(file)
	var buf = make([][]byte, 5)

	for i := range buf {
		buf[i], err = rdr.ReadLineEndingWithCRLF()
		if err != nil {
			return stn, err
		}
	}

	// Server Host & Port.
	stn.ServerHost = strings.TrimSpace(string(buf[0]))

	stn.MainPort, err = number.ParseUint16(strings.TrimSpace(string(buf[1])))
	if err != nil {
		return stn, err
	}

	stn.AuxPort, err = number.ParseUint16(strings.TrimSpace(string(buf[2])))
	if err != nil {
		return stn, err
	}

	// Cache Data Settings.
	stn.Data, err = ParseDataSettings(string(buf[3]), string(buf[4]))
	if err != nil {
		return stn, err
	}

	return stn, nil
}

func (stn *Settings) Check() (err error) {
	if len(stn.File) == 0 {
		return errors.New(ce.ErrFileIsNotSet)
	}

	if len(stn.ServerHost) == 0 {
		return errors.New(ce.ErrServerHostIsNotSet)
	}

	if stn.MainPort == 0 {
		return errors.New(ce.ErrServerPortIsNotSet)
	}

	if stn.AuxPort == 0 {
		return errors.New(ce.ErrServerPortIsNotSet)
	}

	err = stn.Data.Check()
	if err != nil {
		return err
	}

	return nil
}
