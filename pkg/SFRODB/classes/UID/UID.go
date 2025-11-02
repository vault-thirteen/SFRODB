package uid

import (
	"fmt"
	"strings"

	"github.com/vault-thirteen/SFRODB/pkg/SFRODB/protocol"
)

const (
	UidForbiddenChars = "."
)

const (
	ErrNotValid = "UID is not valid: %s"
)

type UID string

func New(s string) (u *UID, err error) {
	u = new(UID)
	*u = UID(strings.TrimSpace(s))

	if !u.isValid() {
		return nil, fmt.Errorf(ErrNotValid, s)
	}

	return u, nil
}

func (u UID) isValid() (ok bool) {
	if strings.ContainsAny(string(u), UidForbiddenChars) {
		return false
	}

	if len(u) > protocol.UidLenMax {
		return false
	}

	return true
}

func (u UID) Length() int {
	return len(u)
}

func (u UID) Bytes() []byte {
	return []byte(u)
}

func (u UID) String() string { return string(u) }
