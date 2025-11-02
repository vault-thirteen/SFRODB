package protocol

import (
	"github.com/vault-thirteen/SFRODB/pkg/SFRODB/classes/Endianness"
)

const (
	LowLevelProtocol      = "tcp"
	TcpKeepAliveIsEnabled = true
	TcpKeepAlivePeriodSec = 15
	Endianness            = endianness.Endianness_BigEndian

	RequestSizeLen  = 2
	ResponseSizeLen = 4
	MethodNameLen   = 3
	StatusNameLen   = 3
	UidLenMax       = 255
	ContentLenMax   = 4_294_967_295 - StatusNameLen
)

// Method strings.
const (
	Method_CloseConnection = "CCC"
	Method_ShowData        = "CSD"
	Method_SearchRecord    = "CSR"
	Method_SearchFile      = "CSF"
	Method_ForgetRecord    = "CFR"
	Method_ResetCache      = "CRC"
)

// Status strings.
const (
	Status_OK                 = "SOK"
	Status_ClientError        = "SER"
	Status_ClosingConnection  = "SCC"
	Status_ShowingData        = "SSD"
	Status_RecordExists       = "SRE"
	Status_RecordDoesNotExist = "SRN"
	Status_FileExists         = "SFE"
	Status_FileDoesNotExist   = "SFN"
)
