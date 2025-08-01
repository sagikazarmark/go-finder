package finder

import (
	"errors"
	"io/fs"
)

// Exists checks if a given file or directory exists.
func Exists(fsys fs.FS, path string, fileType FileType) (bool, error) {
	stat, err := fs.Stat(fsys, path)
	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return fileType.match(stat), nil
}
