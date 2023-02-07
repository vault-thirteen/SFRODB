package method

import (
	"github.com/vault-thirteen/SFRODB/common/method/name"
)

// Method values.
const (
	ClientError       = Method(1)
	OK                = Method(2)
	CloseConnection   = Method(4 + 0)
	ClosingConnection = Method(4 + 2 + 1)

	ShowText      = Method(8 + 1)
	ShowBinary    = Method(8 + 2)
	ShowingText   = Method(8 + 4 + 1)
	ShowingBinary = Method(8 + 4 + 2)

	SearchTextRecord         = Method(16 + 1)
	SearchBinaryRecord       = Method(16 + 2)
	TextRecordExists         = Method(16 + 4 + 1)
	BinaryRecordExists       = Method(16 + 4 + 2)
	TextRecordDoesNotExist   = Method(16 + 8 + 1)
	BinaryRecordDoesNotExist = Method(16 + 8 + 2)

	SearchTextFile         = Method(32 + 1)
	SearchBinaryFile       = Method(32 + 2)
	TextFileExists         = Method(32 + 4 + 1)
	BinaryFileExists       = Method(32 + 4 + 2)
	TextFileDoesNotExist   = Method(32 + 8 + 1)
	BinaryFileDoesNotExist = Method(32 + 8 + 2)

	ForgetTextRecord   = Method(64 + 1)
	ForgetBinaryRecord = Method(64 + 2)
	ResetTextCache     = Method(64 + 32 + 16 + 1)
	ResetBinaryCache   = Method(64 + 32 + 16 + 2)
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
		mn.ShowText,
		mn.ShowBinary,
		mn.ShowingText,
		mn.ShowingBinary,
		mn.SearchTextRecord,
		mn.SearchBinaryRecord,
		mn.TextRecordExists,
		mn.BinaryRecordExists,
		mn.TextRecordDoesNotExist,
		mn.BinaryRecordDoesNotExist,
		mn.SearchTextFile,
		mn.SearchBinaryFile,
		mn.TextFileExists,
		mn.BinaryFileExists,
		mn.TextFileDoesNotExist,
		mn.BinaryFileDoesNotExist,
		mn.ForgetTextRecord,
		mn.ForgetBinaryRecord,
		mn.ResetTextCache,
		mn.ResetBinaryCache,
	}
}

func MethodValues() []Method {
	return []Method{
		ClientError,
		OK,
		CloseConnection,
		ClosingConnection,
		ShowText,
		ShowBinary,
		ShowingText,
		ShowingBinary,
		SearchTextRecord,
		SearchBinaryRecord,
		TextRecordExists,
		BinaryRecordExists,
		TextRecordDoesNotExist,
		BinaryRecordDoesNotExist,
		SearchTextFile,
		SearchBinaryFile,
		TextFileExists,
		BinaryFileExists,
		TextFileDoesNotExist,
		BinaryFileDoesNotExist,
		ForgetTextRecord,
		ForgetBinaryRecord,
		ResetTextCache,
		ResetBinaryCache,
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
