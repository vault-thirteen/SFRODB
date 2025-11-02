package method

import (
	"fmt"
	"strings"

	"github.com/vault-thirteen/SFRODB/pkg/SFRODB/protocol"
)

const (
	Method_Unknown         = Method(0)
	Method_CloseConnection = Method(1)
	Method_ShowData        = Method(2)
	Method_SearchRecord    = Method(3)
	Method_SearchFile      = Method(4)
	Method_ForgetRecord    = Method(5)
	Method_ResetCache      = Method(6)
)

const (
	ErrUnknownMethodName = "unknown method name: %s"
	ErrUnsupportedMethod = "unsupported method: %s"
)

type Method byte

func NewFromString(str string) (m Method, err error) {
	methodStr := strings.TrimSpace(str)

	// For a small number of items, if-branching works faster than maps.
	switch methodStr {
	case protocol.Method_ShowData:
		return Method_ShowData, nil

	case protocol.Method_SearchRecord:
		return Method_SearchRecord, nil

	case protocol.Method_CloseConnection:
		return Method_CloseConnection, nil

	case protocol.Method_SearchFile:
		return Method_SearchFile, nil

	case protocol.Method_ForgetRecord:
		return Method_ForgetRecord, nil

	case protocol.Method_ResetCache:
		return Method_ResetCache, nil

	default:
		return Method_Unknown, fmt.Errorf(ErrUnknownMethodName, methodStr)
	}
}

func (m Method) Bytes() ([]byte, error) {
	switch m {
	case Method_CloseConnection:
		return []byte(protocol.Method_CloseConnection), nil

	case Method_ShowData:
		return []byte(protocol.Method_ShowData), nil

	case Method_SearchRecord:
		return []byte(protocol.Method_SearchRecord), nil

	case Method_SearchFile:
		return []byte(protocol.Method_SearchFile), nil

	case Method_ForgetRecord:
		return []byte(protocol.Method_ForgetRecord), nil

	case Method_ResetCache:
		return []byte(protocol.Method_ResetCache), nil

	default:
		return nil, fmt.Errorf(ErrUnknownMethodName, m)
	}
}
