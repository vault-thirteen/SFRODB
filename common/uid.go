package common

import "strings"

const (
	UidForbiddenChars = "."
)

func IsUidValid(uid string) (ok bool) {
	if strings.ContainsAny(uid, UidForbiddenChars) {
		return false
	}

	return true
}
