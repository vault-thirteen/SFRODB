package settings

import (
	"errors"
	"strings"

	ce "github.com/vault-thirteen/SFRODB/common/error"
	"github.com/vault-thirteen/SFRODB/common/helper"
)

const ErrSyntax = "syntax error"

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

	ds.CacheVolumeMax, err = helper.ParseInt(parts[1])
	if err != nil {
		return nil, err
	}

	ds.CachedItemVolumeMax, err = helper.ParseInt(parts[2])
	if err != nil {
		return nil, err
	}

	ds.CachedItemTTL, err = helper.ParseUint(parts[3])
	if err != nil {
		return nil, err
	}

	return ds, nil
}

func (ds *DataSettings) Check() (err error) {
	if len(ds.Folder) == 0 {
		return errors.New(ce.ErrDataFolderIsNotSet)
	}

	if len(ds.FileExtension) == 0 {
		return errors.New(ce.ErrDataFileExtensionIsNotSet)
	}

	if ds.CacheVolumeMax == 0 {
		return errors.New(ce.ErrCacheVolumeMaxIsNotSet)
	}

	if ds.CachedItemVolumeMax == 0 {
		return errors.New(ce.ErrCachedItemVolumeMaxIsNotSet)
	}

	if ds.CachedItemTTL == 0 {
		return errors.New(ce.ErrCachedItemTTLIsNotSet)
	}

	return nil
}
