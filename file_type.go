package finder

import "io/fs"

// FileType represents the kind of entries [Finder] can return.
type FileType int

const (
	FileTypeAll FileType = iota
	FileTypeFile
	FileTypeDir
)

func (ft FileType) matchDirEntry(entry fs.DirEntry) bool {
	switch ft {
	case FileTypeAll:
		return true

	case FileTypeFile:
		return entry.Type().IsRegular()

	case FileTypeDir:
		return entry.Type().IsDir()

	default:
		return false
	}
}

func (ft FileType) matchFileInfo(info fs.FileInfo) bool {
	switch ft {
	case FileTypeAll:
		return true

	case FileTypeFile:
		return info.Mode().IsRegular()

	case FileTypeDir:
		return info.IsDir()

	default:
		return false
	}
}
