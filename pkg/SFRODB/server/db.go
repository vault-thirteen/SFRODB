package server

import (
	"path/filepath"

	"github.com/vault-thirteen/SFRODB/pkg/SFRODB/common"
	ce "github.com/vault-thirteen/SFRODB/pkg/SFRODB/common/error"
)

// getData gets the data either from cache or from file storage.
// Returns a detailed error.
func (srv *Server) getData(uid string, clientId string) (data []byte, cerr *ce.CommonError) {
	// Check the UID.
	if !common.IsUidValid(uid) {
		return nil, ce.NewClientError(ce.ErrUid, 0, clientId)
	}

	// Try to find the data in cache.
	var err error
	data, err = srv.cache.GetRecord(uid)
	if err == nil {
		return data, nil
	}

	// Try the file storage.
	// Add an extension and convert path to the style of a current OS.
	relPath := filepath.Join(uid+srv.settings.Data.FileExtension, "")
	var fileExists bool
	fileExists, data, err = srv.files.GetFileContents(relPath)
	if !fileExists {
		return nil, ce.NewClientError(err.Error(), 0, clientId)
	}
	if err != nil {
		return nil, ce.NewServerError(err.Error(), 0, clientId)
	}

	// Save data in the cache.
	err = srv.cache.AddRecord(uid, data)
	if err != nil {
		return nil, ce.NewServerError(err.Error(), 0, clientId)
	}

	return data, nil
}
