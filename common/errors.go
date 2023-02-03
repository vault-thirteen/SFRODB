package common

const (
	ErrFileIsNotSet                = "file is not set"
	ErrDataFolderIsNotSet          = "data folder is not set"
	ErrDataFileExtensionIsNotSet   = "data file extension is not set"
	ErrCacheVolumeMaxIsNotSet      = "cache volume limit is not set"
	ErrCachedItemVolumeMaxIsNotSet = "cached item volume limit is not set"
	ErrCachedItemTTLIsNotSet       = "cached item TTL is not set"
	ErrServerHostIsNotSet          = "server host is not set"
	ErrServerPortIsNotSet          = "server port is not set"
	ErrClientHostIsNotSet          = "client host is not set"
	ErrClientPortIsNotSet          = "client port is not set"
	ErrResponseMessageLengthLimit  = "response message length limit is not set"
)

const (
	ErrSrsIsNotSupported = "SRS is unsupported: %d"
	ErrSrsReading        = "SRS reading error: "

	ErrRsReading = "RS reading error: "

	ErrReadingMethodAndData   = "error reading method and data: "
	ErrUnsupportedMethodValue = "unsupported method value: %d"
	ErrUnknownMethodName      = "unknown method name: %s"
	ErrMessageIsTooLong       = "message is too long: %d vs %d"
	ErrTextIsTooLong          = "text is too long: %d vs %d"
	ErrUid                    = "uid error"
	ErrUidIsTooLong           = "uid is too long"

	// ErrSomethingWentWrong is an error which a client sees when he waits for
	// the server's reply and something goes wrong. The reason are variable:
	//	1.	The requested resource has an invalid UID;
	//	2.	The requested resource is not available on the server;
	//	3.	An internal server error has occurred.
	ErrSomethingWentWrong = "something went wrong"
)
