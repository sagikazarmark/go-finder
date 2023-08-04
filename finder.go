// Package finder looks for files and directories in an {fs.Fs} filesystem.
package finder

import (
	"errors"
	"io/fs"
	"path"
	"strings"

	"github.com/sourcegraph/conc/pool"
)

// Finder looks for files and directories in an {fs.Fs} filesystem.
type Finder struct {
	// Paths are the locations where {Finder} looks for the results in.
	Paths []string
	Names []string
}

// Find looks for files and directories in an {fs.Fs} filesystem.
func (f Finder) Find(fsys fs.FS) ([]string, error) {
	// Arbitrary go routine limit (TODO: make this a parameter)
	pool := pool.NewWithResults[[]string]().WithMaxGoroutines(5).WithErrors().WithFirstError()

	for _, searchPath := range f.Paths {
		for _, searchName := range f.Names {
			pool.Go(func() ([]string, error) {
				// If the name contains any glob character, perform a glob match
				if strings.ContainsAny(searchName, "*?[]\\^") {
					return globWalkSearch(fsys, searchPath, searchName)
				} else { //nolint
					return statSearch(fsys, searchPath, searchName)
				}
			})
		}
	}

	allResults, err := pool.Wait()
	if err != nil {
		return nil, err
	}

	var results []string

	for _, r := range allResults {
		results = append(results, r...)
	}

	return results, nil
}

func globWalkSearch(fsys fs.FS, searchPath string, searchName string) ([]string, error) {
	var results []string

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

	return results, nil
}

func statSearch(fsys fs.FS, searchPath string, searchName string) ([]string, error) {
	filePath := path.Join(searchPath, searchName)

	_, err := fs.Stat(fsys, filePath)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return []string{filePath}, nil
}
