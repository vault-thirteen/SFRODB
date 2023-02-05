package common

import (
	"fmt"
	"strings"
)

// As opposed to the HTTP protocol, here
// Method is a hybrid of action type and status code.

// Method name settings.
const (
	MethodNameSpacer      = " "
	MethodNameLengthLimit = 3

	MethodNameClientError        = "EEE"
	MethodNameCloseConnection    = "CLC"
	MethodNameClosingConnection  = "BYE"
	MethodNameShowText           = "ST"
	MethodNameShowingText        = "TT"
	MethodNameShowBinary         = "SB"
	MethodNameShowingBinary      = "BB"
	MethodNameForgetTextRecord   = "RRT"
	MethodNameForgetBinaryRecord = "RRB"
	MethodNameResetTextCache     = "RST"
	MethodNameResetBinaryCache   = "RSB"
	MethodNameOK                 = "OK"
)

// Methods.
const (
	MethodClientError        = Method(0)
	MethodCloseConnection    = Method(1)
	MethodClosingConnection  = Method(2)
	MethodShowText           = Method(4)
	MethodShowingText        = Method(8)
	MethodShowBinary         = Method(16)
	MethodShowingBinary      = Method(32)
	MethodForgetTextRecord   = Method(64 + 1)
	MethodForgetBinaryRecord = Method(64 + 2)
	MethodResetTextCache     = Method(128 + 1)
	MethodResetBinaryCache   = Method(128 + 2)
	MethodOK                 = Method(255)
)

type Method byte

func (c *Connection) NewMethodFromBytes(b []byte) (m Method, err error) {
	if len(b) == 3 {
		return c.NewMethodFromString(string(b))
	}

	return c.NewMethodFromString(string(b[0:3]))
}

func (c *Connection) NewMethodFromString(s string) (m Method, err error) {
	methodStr := strings.TrimSuffix(s, MethodNameSpacer)

	var ok bool
	m, ok = (*c.methodValues)[methodStr]
	if !ok {
		return 0, fmt.Errorf(ErrUnknownMethodName, methodStr)
	}

	return m, nil
}

func MethodNames() []string {
	return []string{
		MethodNameClientError,
		MethodNameCloseConnection,
		MethodNameClosingConnection,
		MethodNameShowText,
		MethodNameShowingText,
		MethodNameShowBinary,
		MethodNameShowingBinary,
		MethodNameForgetTextRecord,
		MethodNameForgetBinaryRecord,
		MethodNameResetTextCache,
		MethodNameResetBinaryCache,
		MethodNameOK,
	}
}

func MethodValues() []Method {
	return []Method{
		MethodClientError,
		MethodCloseConnection,
		MethodClosingConnection,
		MethodShowText,
		MethodShowingText,
		MethodShowBinary,
		MethodShowingBinary,
		MethodForgetTextRecord,
		MethodForgetBinaryRecord,
		MethodResetTextCache,
		MethodResetBinaryCache,
		MethodOK,
	}
}

func InitMethods() (methodNameBuffersMap map[Method][]byte, methodValuesMap map[string]Method) {
	methodNames := MethodNames()
	methodValues := MethodValues()

	methodNameBuffersMap = initMethodNames(methodNames, methodValues)
	methodValuesMap = initMethodValues(methodNames, methodValues)

	return methodNameBuffersMap, methodValuesMap
}

func initMethodNames(methodNames []string, methodValues []Method) (methodNameBuffersMap map[Method][]byte) {
	methodNameBuffersMap = make(map[Method][]byte)

	for i, methodValue := range methodValues {
		methodName := methodNames[i]

		// Create an empty buffer.
		buf := make([]byte, 3)
		for j := range buf {
			buf[j] = MethodNameSpacer[0]
		}

		// Write the method name into the buffer.
		buf[0] = methodName[0]
		if len(methodName) >= 2 {
			buf[1] = methodName[1]
		}
		if len(methodName) >= 3 {
			buf[2] = methodName[2]
		}

		// Save the buffer into the map.
		methodNameBuffersMap[methodValue] = buf
	}

	return methodNameBuffersMap
}

func initMethodValues(methodNames []string, methodValues []Method) (methodValuesMap map[string]Method) {
	methodValuesMap = make(map[string]Method)

	for i, methodName := range methodNames {
		methodValue := methodValues[i]
		methodValuesMap[methodName] = methodValue
	}

	return methodValuesMap
}
