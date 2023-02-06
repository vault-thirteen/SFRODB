package common

import (
	"fmt"
	"strings"
)

// Method name settings.
const (
	// MethodNameSpacer is used to fill free space when method name is shorter
	// than its maximum length (3).
	MethodNameSpacer = " "

	MethodNameLengthLimit = 3
)

// Method names.
const (
	MethodNameClientError       = "EEE"
	MethodNameOK                = "KKK"
	MethodNameCloseConnection   = "CCC"
	MethodNameClosingConnection = "YYY"

	MethodNameShowText      = "ST"
	MethodNameShowBinary    = "SB"
	MethodNameShowingText   = "GT"
	MethodNameShowingBinary = "GB"

	MethodNameSearchTextRecord         = "FRT"
	MethodNameSearchBinaryRecord       = "FRB"
	MethodNameTextRecordExists         = "RET"
	MethodNameBinaryRecordExists       = "REB"
	MethodNameTextRecordDoesNotExist   = "RNT"
	MethodNameBinaryRecordDoesNotExist = "RNB"

	MethodNameSearchTextFile         = "FFT"
	MethodNameSearchBinaryFile       = "FFB"
	MethodNameTextFileExists         = "FET"
	MethodNameBinaryFileExists       = "FEB"
	MethodNameTextFileDoesNotExist   = "FNT"
	MethodNameBinaryFileDoesNotExist = "FNB"

	MethodNameForgetTextRecord   = "RRT"
	MethodNameForgetBinaryRecord = "RRB"
	MethodNameResetTextCache     = "RST"
	MethodNameResetBinaryCache   = "RSB"
)

// Method values.
const (
	MethodClientError       = Method(1)
	MethodOK                = Method(2)
	MethodCloseConnection   = Method(4 + 0)
	MethodClosingConnection = Method(4 + 2 + 1)

	MethodShowText      = Method(8 + 1)
	MethodShowBinary    = Method(8 + 2)
	MethodShowingText   = Method(8 + 4 + 1)
	MethodShowingBinary = Method(8 + 4 + 2)

	MethodSearchTextRecord         = Method(16 + 1)     //TODO
	MethodSearchBinaryRecord       = Method(16 + 2)     //TODO
	MethodTextRecordExists         = Method(16 + 4 + 1) //TODO
	MethodBinaryRecordExists       = Method(16 + 4 + 2) //TODO
	MethodTextRecordDoesNotExist   = Method(16 + 8 + 1) //TODO
	MethodBinaryRecordDoesNotExist = Method(16 + 8 + 2) //TODO

	MethodSearchTextFile         = Method(32 + 1)     //TODO
	MethodSearchBinaryFile       = Method(32 + 2)     //TODO
	MethodTextFileExists         = Method(32 + 4 + 1) //TODO
	MethodBinaryFileExists       = Method(32 + 4 + 2) //TODO
	MethodTextFileDoesNotExist   = Method(32 + 8 + 1) //TODO
	MethodBinaryFileDoesNotExist = Method(32 + 8 + 2) //TODO

	MethodForgetTextRecord   = Method(64 + 1)
	MethodForgetBinaryRecord = Method(64 + 2)
	MethodResetTextCache     = Method(64 + 32 + 16 + 1)
	MethodResetBinaryCache   = Method(64 + 32 + 16 + 2)
)

// Method is a hybrid of action type and status code, as opposed to the HTTP
// protocol.
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
		MethodNameOK,
		MethodNameCloseConnection,
		MethodNameClosingConnection,
		MethodNameShowText,
		MethodNameShowBinary,
		MethodNameShowingText,
		MethodNameShowingBinary,
		MethodNameSearchTextRecord,
		MethodNameSearchBinaryRecord,
		MethodNameTextRecordExists,
		MethodNameBinaryRecordExists,
		MethodNameTextRecordDoesNotExist,
		MethodNameBinaryRecordDoesNotExist,
		MethodNameSearchTextFile,
		MethodNameSearchBinaryFile,
		MethodNameTextFileExists,
		MethodNameBinaryFileExists,
		MethodNameTextFileDoesNotExist,
		MethodNameBinaryFileDoesNotExist,
		MethodNameForgetTextRecord,
		MethodNameForgetBinaryRecord,
		MethodNameResetTextCache,
		MethodNameResetBinaryCache,
	}
}

func MethodValues() []Method {
	return []Method{
		MethodClientError,
		MethodOK,
		MethodCloseConnection,
		MethodClosingConnection,
		MethodShowText,
		MethodShowBinary,
		MethodShowingText,
		MethodShowingBinary,
		MethodSearchTextRecord,
		MethodSearchBinaryRecord,
		MethodTextRecordExists,
		MethodBinaryRecordExists,
		MethodTextRecordDoesNotExist,
		MethodBinaryRecordDoesNotExist,
		MethodSearchTextFile,
		MethodSearchBinaryFile,
		MethodTextFileExists,
		MethodBinaryFileExists,
		MethodTextFileDoesNotExist,
		MethodBinaryFileDoesNotExist,
		MethodForgetTextRecord,
		MethodForgetBinaryRecord,
		MethodResetTextCache,
		MethodResetBinaryCache,
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
