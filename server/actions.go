package server

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/vault-thirteen/SFRODB/common"
	"github.com/vault-thirteen/SFRODB/common/connection"
	ce "github.com/vault-thirteen/SFRODB/common/error"
	"github.com/vault-thirteen/SFRODB/common/method"
)

// showRecord shows a record.
// Returns a detailed error.
func (srv *Server) showRecord(c *connection.Connection, r *common.Request) (err error) {
	switch r.Method {
	case method.ShowText:
		var text string
		text, err = srv.getText(r.UID)
		if err != nil {
			return err
		}
		return srv.showingText(c, text)

	case method.ShowBinary:
		var data []byte
		data, err = srv.getBinary(r.UID)
		if err != nil {
			return err
		}
		return srv.showingBinary(c, data)

	default:
		return ce.NewServerError(fmt.Sprintf(ce.ErrUnsupportedMethodValue, r.Method), 0)
	}
}

// searchRecord checks existence of a record.
// Returns a detailed error.
func (srv *Server) searchRecord(c *connection.Connection, r *common.Request) (err error) {
	// Check the UID.
	if !common.IsUidValid(r.UID) {
		return ce.NewClientError(ce.ErrUid, 0)
	}

	// Search for the record in cache.
	var recExists bool
	switch r.Method {
	case method.SearchTextRecord:
		recExists = srv.cacheT.RecordExists(r.UID)
		if recExists {
			return srv.textRecordExists(c)
		} else {
			return srv.textRecordDoesNotExist(c)
		}

	case method.SearchBinaryRecord:
		recExists = srv.cacheB.RecordExists(r.UID)
		if recExists {
			return srv.binaryRecordExists(c)
		} else {
			return srv.binaryRecordDoesNotExist(c)
		}

	default:
		return ce.NewServerError(fmt.Sprintf(ce.ErrUnsupportedMethodValue, r.Method), 0)
	}
}

// searchFile checks existence of a file.
// Returns a detailed error.
func (srv *Server) searchFile(c *connection.Connection, r *common.Request) (err error) {
	// Check the UID.
	if !common.IsUidValid(r.UID) {
		return ce.NewClientError(ce.ErrUid, 0)
	}

	// Search for the file in storage.
	var fileExists bool
	var relPath string
	switch r.Method {
	case method.SearchTextFile:
		relPath = filepath.Join(r.UID+srv.settings.TextData.FileExtension, "")
		fileExists, err = srv.filesT.FileExists(relPath)
		if err != nil {
			return ce.NewServerError(err.Error(), 0)
		}
		if fileExists {
			return srv.textFileExists(c)
		} else {
			return srv.textFileDoesNotExist(c)
		}

	case method.SearchBinaryFile:
		relPath = filepath.Join(r.UID+srv.settings.BinaryData.FileExtension, "")
		fileExists, err = srv.filesB.FileExists(relPath)
		if err != nil {
			return ce.NewServerError(err.Error(), 0)
		}
		if fileExists {
			return srv.binaryFileExists(c)
		} else {
			return srv.binaryFileDoesNotExist(c)
		}

	default:
		return ce.NewServerError(fmt.Sprintf(ce.ErrUnsupportedMethodValue, r.Method), 0)
	}
}

// forgetRecord removes a record from cache.
// Returns a detailed error.
func (srv *Server) forgetRecord(c *connection.Connection, r *common.Request) (err error) {
	// Check the UID.
	if !common.IsUidValid(r.UID) {
		return ce.NewClientError(ce.ErrUid, 0)
	}

	// Remove the record from the cache.
	var recExists bool
	switch r.Method {
	case method.ForgetTextRecord:
		recExists, err = srv.cacheT.RemoveRecord(r.UID)
	case method.ForgetBinaryRecord:
		recExists, err = srv.cacheB.RemoveRecord(r.UID)
	default:
		return ce.NewServerError(fmt.Sprintf(ce.ErrUnsupportedMethodValue, r.Method), 0)
	}
	if !recExists {
		return ce.NewClientError(err.Error(), 0)
	}
	if err != nil {
		return ce.NewServerError(err.Error(), 0)
	}

	return srv.ok(c)
}

// resetCache removes all records from cache.
// Returns a detailed error.
func (srv *Server) resetCache(c *connection.Connection, r *common.Request) (err error) {
	log.Println(MsgResettingCache)

	// Clear the cache.
	switch r.Method {
	case method.ResetTextCache:
		err = srv.cacheT.Clear()
	case method.ResetBinaryCache:
		err = srv.cacheB.Clear()
	default:
		return ce.NewServerError(fmt.Sprintf(ce.ErrUnsupportedMethodValue, r.Method), 0)
	}
	if err != nil {
		return ce.NewServerError(err.Error(), 0)
	}

	return srv.ok(c)
}

// processError processes a detailed error.
func (srv *Server) processError(err error) (isServerError bool) {
	detailedError, ok := err.(*ce.CommonError)
	if !ok {
		return false
	}

	if detailedError.IsServerError() {
		log.Println(err)
		return true
	}

	return false
}
