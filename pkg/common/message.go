package common

import (
	"fmt"
)

const (
	CR    = "\r"
	CRLF  = CR + LF
	LF    = "\n"
	Space = " "
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
		"[G] = Get a Record" + Space +
		"[E] = Record Exists ?" + Space +
		"[S] = File Exists ?" + Space +
		"[F] = Forget Cache Record" + Space +
		"[R] = Reset Cache" + Space +
		"[Q] = Quit" + CRLF
	Msg_TestClientHintUid             = "Enter the UID > "
	Msg_TestClientHintDataSize        = "Data is quite large. Do you want to see it ? [Y] = Yes; [N] = No. > "
	Msg_UsingDefaultConfigurationFile = "Using the default configuration file."
	Msg_ClearScreen                   = "\033[2J"
	Msg_MoveCursorTopLeft             = "\033[H"
	Msg_EraseCurrentLine              = "\033[2K\r"

	Msgf_RequestDuration = "Request Duration: %d µs. Data Size: %d Bytes."
)

func FormatRequestDuration(durationMks int64, dataSize int) string {
	return fmt.Sprintf(Msgf_RequestDuration, durationMks, dataSize)
}
