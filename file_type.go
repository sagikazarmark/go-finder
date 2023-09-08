package finder

import "io/fs"

// FileType tells [Finder] to look for a certain type of entry.
type FileType int

const (
	AllFileTypes FileType = iota
	File
	Dir
)

func (ft FileType) matchDirEntry(entry fs.DirEntry) bool {
	switch ft {
	case AllFileTypes:
		return true

	case File:
		return entry.Type().IsRegular()

	case Dir:
		return entry.Type().IsDir()

	default:
		return false
	}
}

func (ft FileType) matchFileInfo(info fs.FileInfo) bool {
	switch ft {
	case AllFileTypes:
		return true

	case File:
		return !info.IsDir()

	case Dir:
		return info.IsDir()

	default:
		return false
	}
}
