package uid

import (
	"strings"

	"github.com/vault-thirteen/auxie/errors"

	"github.com/vault-thirteen/SFRODB/pkg/SFRODB/protocol"
)

const (
	UidForbiddenChars = "."
)

const (
	Errf_UidNotValid = "UID is not valid: %s"
)

type UID string

func New(s string) (u *UID, err error) {
	u = new(UID)
	*u = UID(strings.TrimSpace(s))

	if !u.isValid() {
		return nil, errors.ErrorFromFormattedText(Errf_UidNotValid, s)
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
