package server

import (
	"fmt"
	"log"
	"path/filepath"

	ce "github.com/vault-thirteen/SFRODB/pkg/SFRODB/classes/CommonError"
	"github.com/vault-thirteen/SFRODB/pkg/SFRODB/classes/Connection"
	"github.com/vault-thirteen/SFRODB/pkg/SFRODB/classes/Method"
	"github.com/vault-thirteen/SFRODB/pkg/SFRODB/classes/Request"
)

// act_showData shows a data record.
// Returns a detailed error.
func (srv *Server) act_showData(con *connection.Connection, req *request.Request) (cerr *ce.CommonError) {
	if req.Method != method.Method_ShowData {
		return ce.NewServerError(fmt.Sprintf(method.ErrUnsupportedMethod, req.Method), req.Method, 0, con.ClientId())
	}

	var data []byte
	data, cerr = srv.getData(req.UID.String(), con.ClientId())
	if cerr != nil {
		return cerr
	}

	return srv.respond_showingData(con, data)
}

// act_searchRecord checks existence of a record.
// Returns a detailed error.
func (srv *Server) act_searchRecord(con *connection.Connection, req *request.Request) (cerr *ce.CommonError) {
	if req.Method != method.Method_SearchRecord {
		return ce.NewServerError(fmt.Sprintf(method.ErrUnsupportedMethod, req.Method), req.Method, 0, con.ClientId())
	}

	var recExists = srv.cache.RecordExists(req.UID.String())
	if recExists {
		return srv.respond_recordExists(con)
	} else {
		return srv.respond_recordDoesNotExist(con)
	}
}

// act_searchFile checks existence of a file.
// Returns a detailed error.
func (srv *Server) act_searchFile(con *connection.Connection, req *request.Request) (cerr *ce.CommonError) {
	if req.Method != method.Method_SearchFile {
		return ce.NewServerError(fmt.Sprintf(method.ErrUnsupportedMethod, req.Method), req.Method, 0, con.ClientId())
	}

	var relPath = filepath.Join(req.UID.String()+srv.settings.Data.FileExtension, "")

	fileExists, err := srv.files.FileExists(relPath)
	if err != nil {
		return ce.NewServerError(err.Error(), req.Method, 0, con.ClientId())
	}

	if fileExists {
		return srv.respond_fileExists(con)
	} else {
		return srv.respond_fileDoesNotExist(con)
	}
}

// act_forgetRecord removes a record from cache.
// Returns a detailed error.
func (srv *Server) act_forgetRecord(con *connection.Connection, req *request.Request) (cerr *ce.CommonError) {
	if req.Method != method.Method_ForgetRecord {
		return ce.NewServerError(fmt.Sprintf(method.ErrUnsupportedMethod, req.Method), req.Method, 0, con.ClientId())
	}

	srv.cache.RemoveRecord(req.UID.String())

	return srv.respond_ok(con)
}

// act_resetCache removes all records from cache.
// Returns a detailed error.
func (srv *Server) act_resetCache(con *connection.Connection, req *request.Request) (cerr *ce.CommonError) {
	if req.Method != method.Method_ResetCache {
		return ce.NewServerError(fmt.Sprintf(method.ErrUnsupportedMethod, req.Method), req.Method, 0, con.ClientId())
	}

	log.Println(MsgResettingCache)

	err := srv.cache.Clear()
	if err != nil {
		return ce.NewServerError(err.Error(), req.Method, 0, con.ClientId())
	}

	return srv.respond_ok(con)
}
