package finder

import (
	"testing"
	"testing/fstest"
	"time"
)

func TestExists(t *testing.T) {
	// Create a test filesystem with files and directories
	fsys := fstest.MapFS{
		"file.txt": &fstest.MapFile{
			Data:    []byte("test content"),
			Mode:    0o644,
			ModTime: time.Date(2023, time.August, 4, 21, 5, 0, 0, time.UTC),
		},
		"dir/subfile.txt": &fstest.MapFile{
			Data:    []byte("sub content"),
			Mode:    0o644,
			ModTime: time.Date(2023, time.August, 4, 21, 5, 0, 0, time.UTC),
		},
		"executable": &fstest.MapFile{
			Data:    []byte("#!/bin/bash\necho hello"),
			Mode:    0o755,
			ModTime: time.Date(2023, time.August, 4, 21, 5, 0, 0, time.UTC),
		},
	}

	testCases := []struct {
		name        string
		path        string
		fileType    FileType
		expectedRes bool
		expectedErr bool
	}{
		// File exists tests
		{
			name:        "file exists - FileTypeFile",
			path:        "file.txt",
			fileType:    FileTypeFile,
			expectedRes: true,
			expectedErr: false,
		},
		{
			name:        "file exists - FileTypeAny",
			path:        "file.txt",
			fileType:    FileTypeAny,
			expectedRes: true,
			expectedErr: false,
		},
		{
			name:        "file exists - FileTypeAll (deprecated)",
			path:        "file.txt",
			fileType:    FileTypeAll,
			expectedRes: true,
			expectedErr: false,
		},
		{
			name:        "file exists but checking for dir - FileTypeDir",
			path:        "file.txt",
			fileType:    FileTypeDir,
			expectedRes: false,
			expectedErr: false,
		},
		{
			name:        "executable file exists - FileTypeFile",
			path:        "executable",
			fileType:    FileTypeFile,
			expectedRes: true,
			expectedErr: false,
		},

		// Directory exists tests
		{
			name:        "directory exists - FileTypeDir",
			path:        "dir",
			fileType:    FileTypeDir,
			expectedRes: true,
			expectedErr: false,
		},
		{
			name:        "directory exists - FileTypeAny",
			path:        "dir",
			fileType:    FileTypeAny,
			expectedRes: true,
			expectedErr: false,
		},
		{
			name:        "directory exists but checking for file - FileTypeFile",
			path:        "dir",
			fileType:    FileTypeFile,
			expectedRes: false,
			expectedErr: false,
		},
		{
			name:        "subdirectory file exists - FileTypeFile",
			path:        "dir/subfile.txt",
			fileType:    FileTypeFile,
			expectedRes: true,
			expectedErr: false,
		},

		// File/directory does not exist tests
		{
			name:        "file does not exist - FileTypeFile",
			path:        "nonexistent.txt",
			fileType:    FileTypeFile,
			expectedRes: false,
			expectedErr: false,
		},
		{
			name:        "file does not exist - FileTypeDir",
			path:        "nonexistent.txt",
			fileType:    FileTypeDir,
			expectedRes: false,
			expectedErr: false,
		},
		{
			name:        "file does not exist - FileTypeAny",
			path:        "nonexistent.txt",
			fileType:    FileTypeAny,
			expectedRes: false,
			expectedErr: false,
		},
		{
			name:        "directory does not exist - FileTypeDir",
			path:        "nonexistent_dir",
			fileType:    FileTypeDir,
			expectedRes: false,
			expectedErr: false,
		},
		{
			name:        "nested path does not exist",
			path:        "nonexistent/dir/file.txt",
			fileType:    FileTypeFile,
			expectedRes: false,
			expectedErr: false,
		},

		// Edge cases
		{
			name:        "empty path - FileTypeAny",
			path:        "",
			fileType:    FileTypeAny,
			expectedRes: false,
			expectedErr: false,
		},
		{
			name:        "root path - FileTypeDir",
			path:        ".",
			fileType:    FileTypeDir,
			expectedRes: true,
			expectedErr: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			exists, err := Exists(fsys, testCase.path, testCase.fileType)

			if testCase.expectedErr && err == nil {
				t.Errorf("expected an error but got none")
			}
			if !testCase.expectedErr && err != nil {
				t.Errorf("expected no error but got: %v", err)
			}
			if exists != testCase.expectedRes {
				t.Errorf("expected exists to be %v but got %v", testCase.expectedRes, exists)
			}
		})
	}
}
