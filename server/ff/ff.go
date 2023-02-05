package ff

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/vault-thirteen/errorz"
	"github.com/vault-thirteen/file"
)

const (
	ErrFolderIsNotFound = "folder is not found: %s"
	ErrFileDoesNotExist = "file does not exist: "
)

type FilesFolder struct {
	folder string
}

func NewFilesFolder(baseFolder string) (ff *FilesFolder, err error) {
	var ok bool
	ok, err = DoesFolderExist(baseFolder)
	if err != nil {
		return nil, err
	}
	if !ok {
		err = fmt.Errorf(ErrFolderIsNotFound, baseFolder)
		return nil, err
	}

	ff = &FilesFolder{
		folder: baseFolder,
	}

	return ff, nil
}

func (ff *FilesFolder) GetFileContents(relPath string) (fileExists bool, data []byte, err error) {
	filePath := filepath.Join(ff.folder, relPath)

	fileExists, err = file.Exists(filePath)
	if !fileExists {
		return fileExists, nil, errors.New(ErrFileDoesNotExist + filePath)
	}

	var f *os.File
	f, err = os.Open(filePath)
	if err != nil {
		return fileExists, nil, err
	}

	defer func() {
		derr := f.Close()
		if derr != nil {
			err = errorz.Combine(err, derr)
		}
	}()

	data, err = io.ReadAll(f)
	if err != nil {
		return fileExists, nil, err
	}

	return fileExists, data, nil
}

func DoesFolderExist(folderPath string) (ok bool, err error) {
	var info os.FileInfo
	info, err = os.Stat(folderPath)

	if err == nil && info.IsDir() {
		return true, nil
	}

	return false, err
}
