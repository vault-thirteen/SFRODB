package common

import (
	"io"
	"net"
	"time"

	"github.com/vault-thirteen/auxie/cio"
)

const (
	HostPortFormat = "%s:%d"
)

const (
	HttpHeaderFormat_CacheControl = "max-age=%d, must-revalidate"
)

func FormatHostPort(host string, port uint16) string {
	return cio.FormatMessage(HostPortFormat, host, port)
}

func FormatHttpHeader_CacheControl(httpCacheControlMaxAge uint) string {
	return cio.FormatMessage(HttpHeaderFormat_CacheControl, httpCacheControlMaxAge)
}

func EnableTcpKeepAlives(conn *net.TCPConn, kaState bool, periodSec int) (err error) {
	err = conn.SetKeepAlivePeriod(time.Second * time.Duration(periodSec))
	if err != nil {
		return err
	}

	err = conn.SetKeepAlive(kaState)
	if err != nil {
		return err
	}

	return nil
}

func TcpReadExactSize(tc *net.TCPConn, bytesCountToRead uint) (data []byte, err error) {
	data = make([]byte, bytesCountToRead)

	_, err = io.ReadFull(tc, data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
