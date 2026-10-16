package common

import (
	"fmt"
	"log"
)

const (
	CR   = "\r"
	CRLF = CR + LF
	LF   = "\n"
)

const (
	App_Server     = "Server"
	App_TestClient = "Test Client"
)

const (
	Msg_ConnectedTo              = "Connected to "
	Msg_DbClientA                = "DB Client A: "
	Msg_DbClientB                = "DB Client B: "
	Msg_DbErrorListenerStopped   = "DB error listener has stopped."
	Msg_FileExists               = "File exists."
	Msg_FileDoesNotExist         = "File does not exist."
	Msg_HttpErrorListenerStopped = "HTTP error listener has stopped."
	Msg_HttpServer               = "HTTP Server: "
	Msg_HttpsServer              = "HTTPS Server: "
	Msg_HorizontalLine           = "------------------------------------------------------------"
	Msg_Item                     = "Item:"
	Msg_ListenerMain             = "Main Listener: "
	Msg_ListenerAuxiliary        = "Auxiliary Listener: "
	Msg_OsQuitSignal             = "quit signal from OS has been received: "
	Msg_RecordExists             = "Record exists."
	Msg_RecordDoesNotExist       = "Record does not exist."
	Msg_Settings                 = "Settings:"
	Msg_ServerStarting           = "Server is starting ..."
	Msg_ServerStopping           = "Stopping the server ..."
	Msg_ServerStopped            = "Server was stopped."
	Msg_TestClientHintMain       = CRLF +
		"[G] = Get/Show a Record;" + CRLF +
		"[E] = Check Record's Existence;" + CRLF +
		"[S] = Check File's Existence;" + CRLF +
		"[F] = Forget/Remove a Record from Cache;" + CRLF +
		"[R] = Reset/Clear the Cache;" + CRLF +
		"[Q] = Quit/Exit." + CRLF +
		"> "
	Msg_TestClientHintUid             = "Enter the UID > "
	Msg_TestClientHintDataSize        = "Data is quite large. Do you want to see it ? [Y] = Yes; [N] = No. > "
	Msg_UsingDefaultConfigurationFile = "Using the default configuration file."

	Msgf_RequestDuration = "Request Duration: %d µs. Data Size: %d Bytes."
)

func FormatMessage(format string, args ...any) string {
	return fmt.Sprintf(format, args...)
}

func FormatRequestDuration(durationMks int64, dataSize int) string {
	return fmt.Sprintf(Msgf_RequestDuration, durationMks, dataSize)
}

func LogMessage(args ...any) {
	log.Println(args...)
}

func PrintMessage(args ...any) {
	fmt.Println(args...)
}
