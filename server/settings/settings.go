package settings

import (
	"errors"
	"os"
	"strings"

	"github.com/vault-thirteen/SFRODB/common"
	"github.com/vault-thirteen/errorz"
	"github.com/vault-thirteen/reader"
)

type Settings struct {
	// Path to the settings file.
	// Settings are positional for simplicity.
	File string

	// This server's host name.
	ServerHost string

	// This server's port.
	ServerPort uint16

	// Textual and binary data use separate caches.
	TextData   *DataSettings
	BinaryData *DataSettings
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
	var buf = make([][]byte, 6)

	for i := range buf {
		buf[i], err = rdr.ReadLineEndingWithCRLF()
		if err != nil {
			return stn, err
		}
	}

	// Server Host & Port.
	stn.ServerHost = strings.TrimSpace(string(buf[0]))
	stn.ServerPort, err = common.ParseUint16(strings.TrimSpace(string(buf[1])))
	if err != nil {
		return stn, err
	}

	stn.TextData, err = ParseDataSettings(string(buf[2]), string(buf[3]))
	if err != nil {
		return stn, err
	}

	stn.BinaryData, err = ParseDataSettings(string(buf[4]), string(buf[5]))
	if err != nil {
		return stn, err
	}

	return stn, nil
}

func (stn *Settings) Check() (err error) {
	if len(stn.File) == 0 {
		return errors.New(common.ErrFileIsNotSet)
	}

	if len(stn.ServerHost) == 0 {
		return errors.New(common.ErrServerHostIsNotSet)
	}

	if stn.ServerPort == 0 {
		return errors.New(common.ErrServerPortIsNotSet)
	}

	err = stn.TextData.Check()
	if err != nil {
		return err
	}

	err = stn.BinaryData.Check()
	if err != nil {
		return err
	}

	return nil
}
