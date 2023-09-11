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

	// Output: [etc/config.yaml home/user/config.yaml]
}

func TestFinder_Find(t *testing.T) {
	fsys := fstest.MapFS{
		"home/user/.config/app/config.yaml": &fstest.MapFile{
			Mode:    0777,
			ModTime: time.Date(2023, time.August, 4, 21, 5, 0, 0, time.UTC),
		},
		"home/user/app/config.yaml": &fstest.MapFile{
			Mode:    0777,
			ModTime: time.Date(2023, time.August, 4, 21, 5, 0, 0, time.UTC),
		},
		"home/user/config.json": &fstest.MapFile{
			Mode:    0777,
			ModTime: time.Date(2023, time.August, 4, 21, 5, 0, 0, time.UTC),
		},
		"home/user/config.yaml": &fstest.MapFile{
			Mode:    0777,
			ModTime: time.Date(2023, time.August, 4, 21, 5, 0, 0, time.UTC),
		},
		"home/user/config/app.yaml": &fstest.MapFile{
			Mode:    0777,
			ModTime: time.Date(2023, time.August, 4, 21, 5, 0, 0, time.UTC),
		},
		"home/user/config/config.yaml": &fstest.MapFile{
			Mode:    0777,
			ModTime: time.Date(2023, time.August, 4, 21, 5, 0, 0, time.UTC),
		},
		"etc/app/config.yaml": &fstest.MapFile{
			Mode:    0777,
			ModTime: time.Date(2023, time.August, 4, 21, 5, 0, 0, time.UTC),
		},
		"etc/config.json": &fstest.MapFile{
			Mode:    0777,
			ModTime: time.Date(2023, time.August, 4, 21, 5, 0, 0, time.UTC),
		},
		"etc/config.yaml": &fstest.MapFile{
			Mode:    0777,
			ModTime: time.Date(2023, time.August, 4, 21, 5, 0, 0, time.UTC),
		},
		"etc/config/app.yaml": &fstest.MapFile{
			Mode:    0777,
			ModTime: time.Date(2023, time.August, 4, 21, 5, 0, 0, time.UTC),
		},
		"etc/config/config.yaml": &fstest.MapFile{
			Mode:    0777,
			ModTime: time.Date(2023, time.August, 4, 21, 5, 0, 0, time.UTC),
		},

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
			name: "no names to find",
			finder: Finder{
				Paths: []string{"home/user"},
			},
			results: nil,
		},
		{
			name: "no paths to find in",
			finder: Finder{
				Names: []string{"config.yaml"},
			},
			results: nil,
		},
		{
			name: "find in path",
			finder: Finder{
				Paths: []string{"home/user"},
				Names: []string{"config.yaml"},
			},
			results: []string{
				"home/user/config.yaml",
			},
		},
		{
			name: "find in multiple paths",
			finder: Finder{
				Paths: []string{"home/user", "etc"},
				Names: []string{"config.yaml"},
			},
			results: []string{
				"etc/config.yaml",
				"home/user/config.yaml",
			},
		},
		{
			name: "find multiple names in multiple paths",
			finder: Finder{
				Paths: []string{"home/user", "etc"},
				Names: []string{"config", "config.yaml"},
			},
			results: []string{
				"etc/config",
				"etc/config.yaml",
				"home/user/config",
				"home/user/config.yaml",
			},
		},
		{
			name: "find in subdirs of each other",
			finder: Finder{
				Paths: []string{"home/user", "home/user/app"},
				Names: []string{"config.yaml"},
			},
			results: []string{
				"home/user/app/config.yaml",
				"home/user/config.yaml",
			},
		},
		{
			name: "find files only",
			finder: Finder{
				Paths: []string{"home/user", "etc"},
				Names: []string{"config", "config.yaml"},
				Type:  FileTypeFile,
			},
			results: []string{
				"etc/config.yaml",
				"home/user/config.yaml",
			},
		},
		{
			name: "find dirs only",
			finder: Finder{
				Paths: []string{"home/user", "etc"},
				Names: []string{"config", "config.yaml"},
				Type:  FileTypeDir,
			},
			results: []string{
				"etc/config",
				"home/user/config",
			},
		},
		{
			name: "glob match",
			finder: Finder{
				Paths: []string{"home/user", "etc"},
				Names: []string{"config*"},
			},
			results: []string{
				"etc/config",
				"etc/config.json",
				"etc/config.yaml",
				"home/user/config",
				"home/user/config.json",
				"home/user/config.yaml",
			},
		},
		{
			name: "glob match",
			finder: Finder{
				Paths: []string{"home/user", "etc"},
				Names: []string{"config.*"},
			},
			results: []string{
				"etc/config.json",
				"etc/config.yaml",
				"home/user/config.json",
				"home/user/config.yaml",
			},
		},
		{
			name: "glob match files",
			finder: Finder{
				Paths: []string{"home/user", "etc"},
				Names: []string{"config*"},
				Type:  FileTypeFile,
			},
			results: []string{
				"etc/config.json",
				"etc/config.yaml",
				"home/user/config.json",
				"home/user/config.yaml",
			},
		},
		{
			name: "glob match dirs",
			finder: Finder{
				Paths: []string{"home/user", "etc"},
				Names: []string{"config*"},
				Type:  FileTypeDir,
			},
			results: []string{
				"etc/config",
				"home/user/config",
			},
		},
		{
			name: "glob match in subdirs of each other",
			finder: Finder{
				Paths: []string{"home/user", "home/user/config", "etc", "etc/config"},
				Names: []string{"config*"},
			},
			results: []string{
				"etc/config",
				"etc/config.json",
				"etc/config.yaml",
				"etc/config/config.yaml",
				"home/user/config",
				"home/user/config.json",
				"home/user/config.yaml",
				"home/user/config/config.yaml",
			},
		},
		{
			name: "glob match files in subdirs of each other",
			finder: Finder{
				Paths: []string{"home/user", "home/user/config", "etc", "etc/config"},
				Names: []string{"config*"},
				Type:  FileTypeFile,
			},
			results: []string{
				"etc/config.json",
				"etc/config.yaml",
				"etc/config/config.yaml",
				"home/user/config.json",
				"home/user/config.yaml",
				"home/user/config/config.yaml",
			},
		},
		{
			name: "glob match dirs in subdirs of each other",
			finder: Finder{
				Paths: []string{"home/user", "home/user/config", "etc", "etc/config"},
				Names: []string{"config*"},
				Type:  FileTypeDir,
			},
			results: []string{
				"etc/config",
				"home/user/config",
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
