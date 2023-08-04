// Package finder looks for files and directories in an {fs.Fs} filesystem.
package finder

import (
	"errors"
	"io/fs"
	"path"
	"strings"
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
			// If the name contains any glob character, perform a glob match
			if strings.ContainsAny(searchName, "*?[]\\^") {
				err := fs.WalkDir(fsys, searchPath, func(p string, d fs.DirEntry, err error) error {
					if err != nil {
						return err
					}

					// Skip the root
					if p == searchPath {
						return nil
					}

					match, err := path.Match(searchName, d.Name())
					if err != nil {
						return err
					}

					if match {
						results = append(results, p)
					}

					return nil
				})
				if err != nil {
					return results, err
				}
			} else {
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
	}

	return results, nil
}
