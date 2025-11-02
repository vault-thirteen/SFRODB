package tcp

import (
	"net"
	"time"
)

func EnableKeepAlives(conn *net.TCPConn, kaState bool, periodSec int) (err error) {
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
