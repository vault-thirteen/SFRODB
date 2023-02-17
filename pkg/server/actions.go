package server

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/vault-thirteen/SFRODB/pkg/common"
	"github.com/vault-thirteen/SFRODB/pkg/common/connection"
	"github.com/vault-thirteen/SFRODB/pkg/common/error"
	"github.com/vault-thirteen/SFRODB/pkg/common/method"
	"github.com/vault-thirteen/SFRODB/pkg/common/request"
)

// showRecord shows a record.
// Returns a detailed error.
func (srv *Server) showRecord(con *connection.Connection, r *request.Request) (cerr *ce.CommonError) {
	switch r.Method {
	case method.ShowText:
		var text string
		text, cerr = srv.getText(r.UID, con.GetClientId())
		if cerr != nil {
			return cerr
		}
		return srv.showingText(con, text)

	case method.ShowBinary:
		var data []byte
		data, cerr = srv.getBinary(r.UID, con.GetClientId())
		if cerr != nil {
			return cerr
		}
		return srv.showingBinary(con, data)

	default:
		return ce.NewServerError(fmt.Sprintf(ce.ErrUnsupportedMethodValue, r.Method), 0, con.GetClientId())
	}
}

// searchRecord checks existence of a record.
// Returns a detailed error.
func (srv *Server) searchRecord(con *connection.Connection, r *request.Request) (cerr *ce.CommonError) {
	// Check the UID.
	if !common.IsUidValid(r.UID) {
		return ce.NewClientError(ce.ErrUid, 0, con.GetClientId())
	}

	// Search for the record in cache.
	var recExists bool
	switch r.Method {
	case method.SearchTextRecord:
		recExists = srv.cacheT.RecordExists(r.UID)
		if recExists {
			return srv.textRecordExists(con)
		} else {
			return srv.textRecordDoesNotExist(con)
		}

	case method.SearchBinaryRecord:
		recExists = srv.cacheB.RecordExists(r.UID)
		if recExists {
			return srv.binaryRecordExists(con)
		} else {
			return srv.binaryRecordDoesNotExist(con)
		}

	default:
		return ce.NewServerError(fmt.Sprintf(ce.ErrUnsupportedMethodValue, r.Method), 0, con.GetClientId())
	}
}

// searchFile checks existence of a file.
// Returns a detailed error.
func (srv *Server) searchFile(con *connection.Connection, r *request.Request) (cerr *ce.CommonError) {
	// Check the UID.
	if !common.IsUidValid(r.UID) {
		return ce.NewClientError(ce.ErrUid, 0, con.GetClientId())
	}

	// Search for the file in storage.
	var relPath string
	switch r.Method {
	case method.SearchTextFile:
		relPath = filepath.Join(r.UID+srv.settings.TextData.FileExtension, "")
		fileExists, err := srv.filesT.FileExists(relPath)
		if err != nil {
			return ce.NewServerError(err.Error(), 0, con.GetClientId())
		}
		if fileExists {
			return srv.textFileExists(con)
		} else {
			return srv.textFileDoesNotExist(con)
		}

	case method.SearchBinaryFile:
		relPath = filepath.Join(r.UID+srv.settings.BinaryData.FileExtension, "")
		fileExists, err := srv.filesB.FileExists(relPath)
		if err != nil {
			return ce.NewServerError(err.Error(), 0, con.GetClientId())
		}
		if fileExists {
			return srv.binaryFileExists(con)
		} else {
			return srv.binaryFileDoesNotExist(con)
		}

	default:
		return ce.NewServerError(fmt.Sprintf(ce.ErrUnsupportedMethodValue, r.Method), 0, con.GetClientId())
	}
}

// forgetRecord removes a record from cache.
// Returns a detailed error.
func (srv *Server) forgetRecord(con *connection.Connection, r *request.Request) (cerr *ce.CommonError) {
	// Check the UID.
	if !common.IsUidValid(r.UID) {
		return ce.NewClientError(ce.ErrUid, 0, con.GetClientId())
	}

	// Remove the record from the cache.
	var recExists bool
	var err error
	switch r.Method {
	case method.ForgetTextRecord:
		recExists, err = srv.cacheT.RemoveRecord(r.UID)
	case method.ForgetBinaryRecord:
		recExists, err = srv.cacheB.RemoveRecord(r.UID)
	default:
		return ce.NewServerError(fmt.Sprintf(ce.ErrUnsupportedMethodValue, r.Method), 0, con.GetClientId())
	}
	if !recExists {
		return ce.NewClientError(err.Error(), 0, con.GetClientId())
	}
	if err != nil {
		return ce.NewServerError(err.Error(), 0, con.GetClientId())
	}

	return srv.ok(con)
}

// resetCache removes all records from cache.
// Returns a detailed error.
func (srv *Server) resetCache(con *connection.Connection, r *request.Request) (cerr *ce.CommonError) {
	log.Println(MsgResettingCache)

	// Clear the cache.
	var err error
	switch r.Method {
	case method.ResetTextCache:
		err = srv.cacheT.Clear()
	case method.ResetBinaryCache:
		err = srv.cacheB.Clear()
	default:
		return ce.NewServerError(fmt.Sprintf(ce.ErrUnsupportedMethodValue, r.Method), 0, con.GetClientId())
	}
	if err != nil {
		return ce.NewServerError(err.Error(), 0, con.GetClientId())
	}

	return srv.ok(con)
}
