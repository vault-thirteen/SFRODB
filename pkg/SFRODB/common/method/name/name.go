package mn

// Method names.
const (
	ClientError       = "SER"
	OK                = "SOK"
	CloseConnection   = "CCC"
	ClosingConnection = "SCC"

	ShowData    = "CSD"
	ShowingData = "SSD"

	SearchRecord       = "CSR"
	RecordExists       = "SRE"
	RecordDoesNotExist = "SRN"

	SearchFile       = "CSF"
	FileExists       = "SFE"
	FileDoesNotExist = "SFN"

	ForgetRecord = "CFR"
	ResetCache   = "CRS"
)

// Method name settings.
const (
	// Spacer is used to fill free space when method name is shorter
	// than its maximum length (3).
	Spacer = " "

	LengthLimit = 3
)
