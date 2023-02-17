package mn

// Method names.
const (
	ClientError       = "EEE"
	OK                = "KKK"
	CloseConnection   = "CCC"
	ClosingConnection = "YYY"

	ShowText      = "ST"
	ShowBinary    = "SB"
	ShowingText   = "GT"
	ShowingBinary = "GB"

	SearchTextRecord         = "FRT"
	SearchBinaryRecord       = "FRB"
	TextRecordExists         = "RET"
	BinaryRecordExists       = "REB"
	TextRecordDoesNotExist   = "RNT"
	BinaryRecordDoesNotExist = "RNB"

	SearchTextFile         = "FFT"
	SearchBinaryFile       = "FFB"
	TextFileExists         = "FET"
	BinaryFileExists       = "FEB"
	TextFileDoesNotExist   = "FNT"
	BinaryFileDoesNotExist = "FNB"

	ForgetTextRecord   = "RRT"
	ForgetBinaryRecord = "RRB"
	ResetTextCache     = "RST"
	ResetBinaryCache   = "RSB"
)

// Method name settings.
const (
	// Spacer is used to fill free space when method name is shorter
	// than its maximum length (3).
	Spacer = " "

	LengthLimit = 3
)
