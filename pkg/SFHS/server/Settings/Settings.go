package Settings

import (
	"os"
	"strings"

	ae "github.com/vault-thirteen/auxie/errors"
	"github.com/vault-thirteen/auxie/number"
	"github.com/vault-thirteen/auxie/reader"

	cs "github.com/vault-thirteen/SFRODB/pkg/SFRODB/classes/ClientSettings"
)

const (
	Err_ServerModeIsNotSet     = "server mode is not set"
	Err_ServerMode             = "server mode error"
	Err_CertFileIsNotSet       = "certificate file is not set"
	Err_KeyFileIsNotSet        = "key file is not set"
	Err_FileExtensionIsNotSet  = "file extension is not set"
	Err_MimeTypeIsNotSet       = "MIME type is not set"
	Err_DbClientPoolSize       = "DB client pool size is not set"
	Err_HttpCacheControlMaxAge = "HTTP cache control max-age error"
	Err_FileIsNotSet           = "file is not set"
	Err_ServerHostIsNotSet     = "server host is not set"
	Err_ServerPortIsNotSet     = "server port is not set"
)

const (
	ContentDisposition       = ContentDispositionInline
	ContentDispositionInline = "inline"
)

const (
	ServerMode_Http  = "HTTP"
	ServerMode_Https = "HTTPS"
)

const (
	ServerModeId_Http  = 1
	ServerModeId_Https = 2
)

// Settings type is Server's Settings.
type Settings struct {
	// Path to the File with these Settings.
	File string

	// Server's Hostname Name.
	ServerHost string

	// Server's Listen Port.
	ServerPort uint16

	// ServerMode is an HTTP mode selector.
	// Possible values are: HTTP and HTTPS.
	ServerModeStr string
	ServerModeId  byte

	// Server's Certificate and Key.
	CertFile string
	KeyFile  string

	// Database Hostname Name.
	DbHost string

	// Database Ports.
	DbPortA uint16
	DbPortB uint16

	// DbClientPoolSize is the size of a pool of DB clients.
	DbClientPoolSize int

	// File Extension & MIME Type.
	// Extension which is appended to all files served.
	FileExtension string
	MimeType      string

	// HttpCacheControlMaxAge is time in seconds for which this server's
	// response is fresh (valid). After this period clients will be refreshing
	// the stale content by re-requesting it from the server.
	HttpCacheControlMaxAge uint

	// Allowed Origin for cross-origin requests (CORS).
	AllowedOriginForCORS string
}

func NewSettingsFromFile(filePath string) (stn *Settings, err error) {
	stn = &Settings{
		File: filePath,
	}

	var file *os.File
	file, err = os.Open(stn.File)
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
	var buf = make([][]byte, 13)

	for i := range buf {
		buf[i], err = rdr.ReadLineEndingWithCRLF()
		if err != nil {
			return stn, err
		}
	}

	// Server Hostname & Port.
	stn.ServerHost = strings.TrimSpace(string(buf[0]))

	stn.ServerPort, err = number.ParseUint16(strings.TrimSpace(string(buf[1])))
	if err != nil {
		return stn, err
	}

	// Server Work Mode.
	stn.ServerModeStr = strings.ToUpper(strings.TrimSpace(string(buf[2])))
	switch stn.ServerModeStr {
	case ServerMode_Http:
		stn.ServerModeId = ServerModeId_Http
	case ServerMode_Https:
		stn.ServerModeId = ServerModeId_Https
	}

	// Certificate and Key for optional TLS.
	stn.CertFile = strings.TrimSpace(string(buf[3]))
	stn.KeyFile = strings.TrimSpace(string(buf[4]))

	// Database.
	stn.DbHost = strings.TrimSpace(string(buf[5]))

	stn.DbPortA, err = number.ParseUint16(strings.TrimSpace(string(buf[6])))
	if err != nil {
		return stn, err
	}

	stn.DbPortB, err = number.ParseUint16(strings.TrimSpace(string(buf[7])))
	if err != nil {
		return stn, err
	}

	stn.DbClientPoolSize, err = number.ParseInt(strings.TrimSpace(string(buf[8])))
	if err != nil {
		return stn, err
	}

	// HTTP.
	stn.FileExtension = strings.TrimSpace(string(buf[9]))
	if stn.FileExtension[0] != '.' {
		stn.FileExtension = `.` + stn.FileExtension
	}

	stn.MimeType = strings.TrimSpace(string(buf[10]))

	stn.HttpCacheControlMaxAge, err = number.ParseUint(strings.TrimSpace(string(buf[11])))
	if err != nil {
		return stn, err
	}

	stn.AllowedOriginForCORS = strings.TrimSpace(string(buf[12]))

	return stn, nil
}

func (stn *Settings) Check() (err error) {
	if len(stn.File) == 0 {
		return ae.ErrorFromText(Err_FileIsNotSet)
	}

	if len(stn.ServerHost) == 0 {
		return ae.ErrorFromText(Err_ServerHostIsNotSet)
	}

	if stn.ServerPort == 0 {
		return ae.ErrorFromText(Err_ServerPortIsNotSet)
	}

	if len(stn.ServerModeStr) == 0 {
		return ae.ErrorFromText(Err_ServerModeIsNotSet)
	} else {
		if (stn.ServerModeStr != ServerMode_Http) &&
			(stn.ServerModeStr != ServerMode_Https) {
			return ae.ErrorFromText(Err_ServerMode)
		}
	}

	if stn.ServerModeId == 0 {
		return ae.ErrorFromText(Err_ServerModeIsNotSet)
	} else {
		if (stn.ServerModeId != ServerModeId_Http) &&
			(stn.ServerModeId != ServerModeId_Https) {
			return ae.ErrorFromText(Err_ServerMode)
		}
	}

	switch stn.ServerModeStr {
	case ServerMode_Http:
		// Keys are not required.
	case ServerMode_Https:
		if len(stn.CertFile) == 0 {
			return ae.ErrorFromText(Err_CertFileIsNotSet)
		}
		if len(stn.KeyFile) == 0 {
			return ae.ErrorFromText(Err_KeyFileIsNotSet)
		}
	default:
		return ae.ErrorFromText(Err_ServerMode)
	}

	if len(stn.DbHost) == 0 {
		return ae.ErrorFromText(cs.Err_ClientHostIsNotSet)
	}

	if stn.DbPortA == 0 {
		return ae.ErrorFromText(cs.Err_ClientPortIsNotSet)
	}

	if stn.DbPortB == 0 {
		return ae.ErrorFromText(cs.Err_ClientPortIsNotSet)
	}

	if stn.DbClientPoolSize == 0 {
		return ae.ErrorFromText(Err_DbClientPoolSize)
	}

	if len(stn.FileExtension) == 0 {
		return ae.ErrorFromText(Err_FileExtensionIsNotSet)
	}

	if len(stn.MimeType) == 0 {
		return ae.ErrorFromText(Err_MimeTypeIsNotSet)
	}

	if stn.HttpCacheControlMaxAge == 0 {
		return ae.ErrorFromText(Err_HttpCacheControlMaxAge)
	}

	// AllowedOriginForCORS is not checked as it may be empty.

	return nil
}
