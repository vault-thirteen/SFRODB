package tcp

import (
	"io"
	"net"
)

func ReadExactSize(tc *net.TCPConn, bytesCountToRead uint) (data []byte, err error) {
	data = make([]byte, bytesCountToRead)

	_, err = io.ReadFull(tc, data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
