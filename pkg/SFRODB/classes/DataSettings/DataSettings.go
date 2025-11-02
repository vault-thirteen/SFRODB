package ds

import (
	"errors"
	"strings"

	"github.com/vault-thirteen/auxie/number"
)

const (
	ErrSyntax                      = "syntax error"
	ErrDataFolderIsNotSet          = "data folder is not set"
	ErrDataFileExtensionIsNotSet   = "data file extension is not set"
	ErrCacheVolumeMaxIsNotSet      = "cache's maximum volume is not set"
	ErrCachedItemVolumeMaxIsNotSet = "cached item's maximum volume is not set"
	ErrCachedItemTTLIsNotSet       = "cached item's TTL is not set"
)

type DataSettings struct {
	// 1. Folder with data files.
	Folder string

	// 2. File extension.
	// This extension is concatenated with UIDs of cached items to get the full
	// file name of the item. Each record (item) is stored in a separate file.
	FileExtension string

	// 3. Maximum size of cache in bytes.
	CacheVolumeMax int

	// 4. Maximum size of a single cached item in bytes.
	CachedItemVolumeMax int

	// 5. Expiration time of a single cached item in seconds.
	// 1 Hour = 3600 Seconds,
	// 1 Day = 86400 Seconds.
	CachedItemTTL uint
}

func ParseDataSettings(line1, line2 string) (ds *DataSettings, err error) {
	ds = &DataSettings{
		Folder: strings.TrimSpace(line1),
	}

	parts := strings.Split(strings.TrimSpace(line2), " ")
	if len(parts) != 4 {
		return nil, errors.New(ErrSyntax)
	}

	ds.FileExtension = parts[0]
	if len(ds.FileExtension) == 0 {
		return nil, errors.New(ErrSyntax)
	}
	if ds.FileExtension[0] != '.' {
		ds.FileExtension = `.` + ds.FileExtension
	}

	ds.CacheVolumeMax, err = number.ParseInt(parts[1])
	if err != nil {
		return nil, err
	}

	ds.CachedItemVolumeMax, err = number.ParseInt(parts[2])
	if err != nil {
		return nil, err
	}

	ds.CachedItemTTL, err = number.ParseUint(parts[3])
	if err != nil {
		return nil, err
	}

	return ds, nil
}

func (ds *DataSettings) Check() (err error) {
	if len(ds.Folder) == 0 {
		return errors.New(ErrDataFolderIsNotSet)
	}

	if len(ds.FileExtension) == 0 {
		return errors.New(ErrDataFileExtensionIsNotSet)
	}

	if ds.CacheVolumeMax == 0 {
		return errors.New(ErrCacheVolumeMaxIsNotSet)
	}

	if ds.CachedItemVolumeMax == 0 {
		return errors.New(ErrCachedItemVolumeMaxIsNotSet)
	}

	if ds.CachedItemTTL == 0 {
		return errors.New(ErrCachedItemTTLIsNotSet)
	}

	return nil
}
