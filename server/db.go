package server

import (
	"github.com/vault-thirteen/SFRODB/common"
	"path/filepath"
)

// getText gets the text either from cache or from file storage.
// Returns a detailed error.
func (srv *Server) getText(uid string) (text string, err error) {
	// Check the UID.
	if !common.IsUidValid(uid) {
		return "", common.NewClientError(common.ErrUid, 0)
	}

	// Try to find the text in cache.
	text, err = srv.cacheT.GetRecord(uid)
	if err == nil {
		return text, nil
	}

	// Try the file storage.
	// Add an extension and convert path to the style of a current OS.
	relPath := filepath.Join(uid+srv.settings.TextData.FileExtension, "")
	var data []byte
	var fileExists bool
	fileExists, data, err = srv.filesT.GetFileContents(relPath)
	if !fileExists {
		return "", common.NewClientError(err.Error(), 0)
	}
	if err != nil {
		return "", common.NewServerError(err.Error(), 0)
	}
	text = string(data)

	// Save data in the cache.
	err = srv.cacheT.AddRecord(uid, text)
	if err != nil {
		return "", common.NewServerError(err.Error(), 0)
	}

	return text, nil
}

// getBinary gets the binary data either from cache or from file storage.
// Returns a detailed error.
func (srv *Server) getBinary(uid string) (data []byte, err error) {
	// Check the UID.
	if !common.IsUidValid(uid) {
		return nil, common.NewClientError(common.ErrUid, 0)
	}

	// Try to find the data in cache.
	data, err = srv.cacheB.GetRecord(uid)
	if err == nil {
		return data, nil
	}

	// Try the file storage.
	// Add an extension and convert path to the style of a current OS.
	relPath := filepath.Join(uid+srv.settings.BinaryData.FileExtension, "")
	var fileExists bool
	fileExists, data, err = srv.filesB.GetFileContents(relPath)
	if !fileExists {
		return nil, common.NewClientError(err.Error(), 0)
	}
	if err != nil {
		return nil, common.NewServerError(err.Error(), 0)
	}

	// Save data in the cache.
	err = srv.cacheB.AddRecord(uid, data)
	if err != nil {
		return nil, common.NewServerError(err.Error(), 0)
	}

	return data, nil
}
