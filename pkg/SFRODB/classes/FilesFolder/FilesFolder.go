package ff

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	ae "github.com/vault-thirteen/auxie/errors"
	"github.com/vault-thirteen/auxie/file"
)

const (
	Err_RelPathIsNotValid = "relative path is not valid"

	Errf_FolderIsNotFound = "folder is not found: %s"
	Errf_FileDoesNotExist = "file does not exist: %s"
)

type FilesFolder struct {
	folder        string
	storageAccess *sync.Mutex
}

func New(baseFolder string) (ff *FilesFolder, err error) {
	var ok bool
	ok, err = file.FolderExists(baseFolder)
	if err != nil {
		return nil, err
	}
	if !ok {
		err = ae.ErrorFromFormattedText(Errf_FolderIsNotFound, baseFolder)
		return nil, err
	}

	ff = &FilesFolder{
		folder:        baseFolder,
		storageAccess: new(sync.Mutex),
	}

	return ff, nil
}

func (ff *FilesFolder) GetFileContents(relPath string) (fileExists bool, data []byte, err error) {
	if !isRelPathValid(relPath) {
		return false, nil, ae.ErrorFromText(Err_RelPathIsNotValid)
	}

	filePath := filepath.Join(ff.folder, relPath)

	ff.storageAccess.Lock()
	defer ff.storageAccess.Unlock()

	fileExists, err = file.FileExists(filePath)
	if !fileExists {
		return false, nil, ae.ErrorFromFormattedText(Errf_FileDoesNotExist, filePath)
	}

	var f *os.File
	f, err = os.Open(filePath)
	if err != nil {
		return fileExists, nil, err
	}

	defer func() {
		derr := f.Close()
		if derr != nil {
			err = ae.Combine(err, derr)
		}
	}()

	data, err = io.ReadAll(f)
	if err != nil {
		return fileExists, nil, err
	}

	return true, data, nil
}

func (ff *FilesFolder) FileExists(relPath string) (fileExists bool, err error) {
	if !isRelPathValid(relPath) {
		return false, ae.ErrorFromText(Err_RelPathIsNotValid)
	}

	filePath := filepath.Join(ff.folder, relPath)

	ff.storageAccess.Lock()
	defer ff.storageAccess.Unlock()

	return file.FileExists(filePath)
}

func isRelPathValid(relPath string) (ok bool) {
	return !strings.Contains(relPath, "..")
}
