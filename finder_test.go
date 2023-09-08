package finder

import (
	"fmt"
	"os"
	"testing"
	"testing/fstest"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Example() {
	fsys := os.DirFS("./testdata/")

	finder := Finder{
		Paths: []string{
			"home/user",
			"etc",
		},
		Names: []string{"config.*"},
		Type:  FileTypeFile,
	}

	results, err := finder.Find(fsys)
	if err != nil {
		panic(err)
	}

	fmt.Print(results)

	// Output: [home/user/config.yaml etc/config.yaml]
}

func TestFinder_Find(t *testing.T) {
	fsys := fstest.MapFS{
		"foo/bar": &fstest.MapFile{
			Data:    []byte("hello world from bar"),
			Mode:    0777,
			ModTime: time.Date(2023, time.August, 4, 21, 5, 0, 0, time.UTC),
		},
		"foo/baz": &fstest.MapFile{
			Data:    []byte("hello world from baz"),
			Mode:    0777,
			ModTime: time.Date(2023, time.August, 4, 21, 5, 0, 0, time.UTC),
		},
		"foo/bat/bar/baz": &fstest.MapFile{
			Data:    []byte("hello world from four levels"),
			Mode:    0777,
			ModTime: time.Date(2023, time.August, 4, 21, 5, 0, 0, time.UTC),
		},
	}

	tests := []struct {
		name    string
		finder  Finder
		results []string
	}{
		{
			name:    "nothing to find",
			finder:  Finder{},
			results: nil,
		},
		{
			name: "find file",
			finder: Finder{
				Paths: []string{"foo"},
				Names: []string{"bar"},
			},
			results: []string{"foo/bar"},
		},
		{
			name: "file does not exist",
			finder: Finder{
				Paths: []string{"foo"},
				Names: []string{"nope"},
			},
			results: nil,
		},
		{
			name: "file is dir",
			finder: Finder{
				Paths: []string{"foo"},
				Names: []string{"bat"},
				Type:  FileTypeFile,
			},
			results: nil,
		},
		{
			name: "dir is file",
			finder: Finder{
				Paths: []string{"foo"},
				Names: []string{"bar"},
				Type:  FileTypeDir,
			},
			results: nil,
		},
		{
			name: "glob match",
			finder: Finder{
				Paths: []string{"foo"},
				Names: []string{"ba?"},
			},
			results: []string{
				"foo/bar",
				"foo/bat",
				"foo/bat/bar",
				"foo/bat/bar/baz",
				"foo/baz",
			},
		},
		{
			name: "glob match files",
			finder: Finder{
				Paths: []string{"foo"},
				Names: []string{"ba?"},
				Type:  FileTypeFile,
			},
			results: []string{
				"foo/bar",
				"foo/bat/bar/baz",
				"foo/baz",
			},
		},

		{
			name: "glob match dirs",
			finder: Finder{
				Paths: []string{"foo"},
				Names: []string{"ba?"},
				Type:  FileTypeDir,
			},
			results: []string{
				"foo/bat",
				"foo/bat/bar",
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			results, err := tt.finder.Find(fsys)
			require.NoError(t, err)

			assert.Equal(t, tt.results, results)
		})
	}
}
