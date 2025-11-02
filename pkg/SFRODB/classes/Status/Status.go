package status

import (
	"fmt"
	"strings"

	"github.com/vault-thirteen/SFRODB/pkg/SFRODB/protocol"
)

const (
	Status_Unknown            = Status(0)
	Status_OK                 = Status(1)
	Status_ClientError        = Status(2)
	Status_ClosingConnection  = Status(3)
	Status_ShowingData        = Status(4)
	Status_RecordExists       = Status(5)
	Status_RecordDoesNotExist = Status(6)
	Status_FileExists         = Status(7)
	Status_FileDoesNotExist   = Status(8)
)

const (
	ErrUnknownStatusName = "unknown status name: %s"
)

type Status byte

func NewFromString(str string) (s Status, err error) {
	statusStr := strings.TrimSpace(str)

	// For a small number of items, if-branching works faster than maps.
	switch statusStr {
	case protocol.Status_OK:
		return Status_OK, nil

	case protocol.Status_ClientError:
		return Status_ClientError, nil

	case protocol.Status_ClosingConnection:
		return Status_ClosingConnection, nil

	case protocol.Status_ShowingData:
		return Status_ShowingData, nil

	case protocol.Status_RecordExists:
		return Status_RecordExists, nil

	case protocol.Status_RecordDoesNotExist:
		return Status_RecordDoesNotExist, nil

	case protocol.Status_FileExists:
		return Status_FileExists, nil

	case protocol.Status_FileDoesNotExist:
		return Status_FileDoesNotExist, nil

	default:
		return Status_Unknown, fmt.Errorf(ErrUnknownStatusName, statusStr)
	}
}

func (s Status) Bytes() ([]byte, error) {
	switch s {
	case Status_OK:
		return []byte(protocol.Status_OK), nil

	case Status_ClientError:
		return []byte(protocol.Status_ClientError), nil

	case Status_ClosingConnection:
		return []byte(protocol.Status_ClosingConnection), nil

	case Status_ShowingData:
		return []byte(protocol.Status_ShowingData), nil

	case Status_RecordExists:
		return []byte(protocol.Status_RecordExists), nil

	case Status_RecordDoesNotExist:
		return []byte(protocol.Status_RecordDoesNotExist), nil

	case Status_FileExists:
		return []byte(protocol.Status_FileExists), nil

	case Status_FileDoesNotExist:
		return []byte(protocol.Status_FileDoesNotExist), nil

	default:
		return nil, fmt.Errorf(ErrUnknownStatusName, s)
	}
}
