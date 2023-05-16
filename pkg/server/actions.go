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

// showData shows a data record.
// Returns a detailed error.
func (srv *Server) showData(con *connection.Connection, r *request.Request) (cerr *ce.CommonError) {
	switch r.Method {
	case method.ShowData:
		var data []byte
		data, cerr = srv.getData(r.UID, con.ClientId())
		if cerr != nil {
			return cerr
		}
		return srv.showingData(con, data)

	default:
		return ce.NewServerError(fmt.Sprintf(ce.ErrUnsupportedMethodValue, r.Method), 0, con.ClientId())
	}
}

// searchRecord checks existence of a record.
// Returns a detailed error.
func (srv *Server) searchRecord(con *connection.Connection, r *request.Request) (cerr *ce.CommonError) {
	// Check the UID.
	if !common.IsUidValid(r.UID) {
		return ce.NewClientError(ce.ErrUid, 0, con.ClientId())
	}

	// Search for the record in cache.
	switch r.Method {
	case method.SearchRecord:
		var recExists = srv.cache.RecordExists(r.UID)
		if recExists {
			return srv.recordExists(con)
		} else {
			return srv.recordDoesNotExist(con)
		}

	default:
		return ce.NewServerError(fmt.Sprintf(ce.ErrUnsupportedMethodValue, r.Method), 0, con.ClientId())
	}
}

// searchFile checks existence of a file.
// Returns a detailed error.
func (srv *Server) searchFile(con *connection.Connection, r *request.Request) (cerr *ce.CommonError) {
	// Check the UID.
	if !common.IsUidValid(r.UID) {
		return ce.NewClientError(ce.ErrUid, 0, con.ClientId())
	}

	// Search for the file in storage.
	switch r.Method {
	case method.SearchFile:
		var relPath = filepath.Join(r.UID+srv.settings.Data.FileExtension, "")
		fileExists, err := srv.files.FileExists(relPath)
		if err != nil {
			return ce.NewServerError(err.Error(), 0, con.ClientId())
		}
		if fileExists {
			return srv.fileExists(con)
		} else {
			return srv.fileDoesNotExist(con)
		}

	default:
		return ce.NewServerError(fmt.Sprintf(ce.ErrUnsupportedMethodValue, r.Method), 0, con.ClientId())
	}
}

// forgetRecord removes a record from cache.
// Returns a detailed error.
func (srv *Server) forgetRecord(con *connection.Connection, r *request.Request) (cerr *ce.CommonError) {
	// Check the UID.
	if !common.IsUidValid(r.UID) {
		return ce.NewClientError(ce.ErrUid, 0, con.ClientId())
	}

	// Remove the record from the cache.
	var recExists bool
	var err error
	switch r.Method {
	case method.ForgetRecord:
		recExists, err = srv.cache.RemoveRecord(r.UID)
	default:
		return ce.NewServerError(fmt.Sprintf(ce.ErrUnsupportedMethodValue, r.Method), 0, con.ClientId())
	}
	if !recExists {
		return ce.NewClientError(err.Error(), 0, con.ClientId())
	}
	if err != nil {
		return ce.NewServerError(err.Error(), 0, con.ClientId())
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
	case method.ResetCache:
		err = srv.cache.Clear()
	default:
		return ce.NewServerError(fmt.Sprintf(ce.ErrUnsupportedMethodValue, r.Method), 0, con.ClientId())
	}
	if err != nil {
		return ce.NewServerError(err.Error(), 0, con.ClientId())
	}

	return srv.ok(con)
}
