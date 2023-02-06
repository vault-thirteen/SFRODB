package ff

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/vault-thirteen/errorz"
	"github.com/vault-thirteen/file"
)

const (
	ErrFolderIsNotFound  = "folder is not found: %s"
	ErrFileDoesNotExist  = "file does not exist: "
	ErrRelPathIsNotValid = "relative path is not valid"
)

type FilesFolder struct {
	folder string
}

func NewFilesFolder(baseFolder string) (ff *FilesFolder, err error) {
	var ok bool
	ok, err = file.FolderExists(baseFolder)
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
	if !isRelPathValid(relPath) {
		return false, nil, errors.New(ErrRelPathIsNotValid)
	}

	filePath := filepath.Join(ff.folder, relPath)

	fileExists, err = file.FileExists(filePath)
	if !fileExists {
		return false, nil, errors.New(ErrFileDoesNotExist + filePath)
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

	return true, data, nil
}

func (ff *FilesFolder) FileExists(relPath string) (fileExists bool, err error) {
	if !isRelPathValid(relPath) {
		return false, errors.New(ErrRelPathIsNotValid)
	}

	filePath := filepath.Join(ff.folder, relPath)

	return file.FileExists(filePath)
}

func isRelPathValid(relPath string) (ok bool) {
	return !strings.Contains(relPath, "..")
}
