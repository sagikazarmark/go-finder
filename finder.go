// Package finder looks for files and directories in an {fs.Fs} filesystem.
package finder

import (
	"errors"
	"io/fs"
	"path"
)

// Finder looks for files and directories in an {fs.Fs} filesystem.
type Finder struct {
	// Paths are the locations where {Finder} looks for the results in.
	Paths []string
	Names []string
}

// Find looks for files and directories in an {fs.Fs} filesystem.
func (f Finder) Find(fsys fs.FS) ([]string, error) {
	var results []string

	for _, searchPath := range f.Paths {
		for _, searchName := range f.Names {
			filePath := path.Join(searchPath, searchName)

			_, err := fs.Stat(fsys, filePath)
			if errors.Is(err, fs.ErrNotExist) {
				continue
			}
			if err != nil {
				return nil, err
			}

			results = append(results, filePath)
		}
	}

	return results, nil
}
