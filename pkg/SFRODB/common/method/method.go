package method

import (
	"github.com/vault-thirteen/SFRODB/pkg/SFRODB/common/method/name"
)

// Method values.
const (
	ClientError = Method(1)
	OK          = Method(255)

	CloseConnection   = Method(2 + 0)
	ClosingConnection = Method(2 + 1)

	ShowData    = Method(4 + 0)
	ShowingData = Method(4 + 2 + 1)

	SearchRecord       = Method(8 + 0)
	RecordExists       = Method(8 + 4 + 2)
	RecordDoesNotExist = Method(8 + 1)

	SearchFile       = Method(16)
	FileExists       = Method(16 + 8 + 4)
	FileDoesNotExist = Method(16 + 2 + 1)

	ForgetRecord = Method(32)

	ResetCache = Method(64)
)

// Method is a hybrid of action type and status code, as opposed to the HTTP
// protocol.
type Method byte

func MethodNames() []string {
	return []string{
		mn.ClientError,
		mn.OK,
		mn.CloseConnection,
		mn.ClosingConnection,
		mn.ShowData,
		mn.ShowingData,
		mn.SearchRecord,
		mn.RecordExists,
		mn.RecordDoesNotExist,
		mn.SearchFile,
		mn.FileExists,
		mn.FileDoesNotExist,
		mn.ForgetRecord,
		mn.ResetCache,
	}
}

func MethodValues() []Method {
	return []Method{
		ClientError,
		OK,
		CloseConnection,
		ClosingConnection,
		ShowData,
		ShowingData,
		SearchRecord,
		RecordExists,
		RecordDoesNotExist,
		SearchFile,
		FileExists,
		FileDoesNotExist,
		ForgetRecord,
		ResetCache,
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
			buf[j] = mn.Spacer[0]
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
