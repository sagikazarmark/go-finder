package finder

import (
	"testing"
	"testing/fstest"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFinder_Find(t *testing.T) {
	fsys := fstest.MapFS{
		"foo/bar": &fstest.MapFile{
			Data:    []byte("hello world"),
			Mode:    0777,
			ModTime: time.Date(2023, time.August, 4, 21, 5, 0, 0, time.UTC),
		},
		"foo/baz": &fstest.MapFile{
			Data:    []byte("hello world"),
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
			name: "file glob match",
			finder: Finder{
				Paths: []string{"foo"},
				Names: []string{"ba?"},
			},
			results: []string{"foo/bar", "foo/baz"},
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
