package tcp

import (
	"errors"
	"net"
)

const (
	ErrUnexpectedExtraData = "unexpected extra data"
)

func ReadExactSize(tc *net.TCPConn, bytesCountToRead uint) (data []byte, err error) {
	data = make([]byte, 0, bytesCountToRead)
	var (
		bytesReceived uint = 0
		buf           []byte
		chunkSize     int
		bytesExpected uint
	)

	for {
		bytesExpected = bytesCountToRead - bytesReceived
		buf = make([]byte, bytesExpected)
		chunkSize, err = tc.Read(buf)
		if err != nil {
			return data, err
		}
		if uint(chunkSize) > bytesExpected {
			return nil, errors.New(ErrUnexpectedExtraData)
		}

		data = append(data, buf[0:chunkSize]...)
		bytesReceived += uint(chunkSize)
		if bytesReceived == bytesCountToRead {
			break
		}
	}

	return data, nil
}
