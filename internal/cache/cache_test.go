package cache

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	store "github.com/yakuninmax/imgpreviewer/internal/storage"
)

func TestCache(t *testing.T) {
	size := int64(500000)

	testFiles := []file{
		{
			uri:  "examples/_gopher_original_1024x504.jpg",
			size: 64212,
			name: "_gopher_original_1024x504",
		},
		{
			uri:  "examples/gopher_50x50.jpg",
			size: 1956,
			name: "gopher_50x50",
		},
		{
			uri:  "examples/gopher_200x700.jpg",
			size: 30146,
			name: "gopher_200x700",
		},
		{
			uri:  "examples/gopher_256x126.jpg",
			size: 10121,
			name: "gopher_256x126",
		},
		{
			uri:  "examples/gopher_333x666.jpg",
			size: 41562,
			name: "gopher_333x666",
		},
		{
			uri:  "examples/gopher_500x500.jpg",
			size: 47656,
			name: "gopher_500x500",
		},
		{
			uri:  "examples/gopher_1024x252.jpg",
			size: 41771,
			name: "gopher_1024x252",
		},
		{
			uri:  "examples/gopher_2000x1000.jpg",
			size: 226943,
			name: "gopher_2000x1000",
		},
	}

	t.Run("put files to cache", func(t *testing.T) {
		s, _ := store.New("/tmp/")
		c, _ := New(size, s)
		for _, file := range testFiles {
			d, _ := os.ReadFile(file.uri)

			err := c.Put(file.uri, d)

			require.Nil(t, err)
		}

		_ = s.Clean()
	})

	t.Run("get file from cache", func(t *testing.T) {
		s, _ := store.New("/tmp/")
		c, _ := New(size, s)

		d, _ := os.ReadFile(testFiles[0].uri)
		err := c.Put(testFiles[0].uri, d)

		require.Nil(t, err)

		cd, err := c.Get(testFiles[0].uri)

		require.Nil(t, err)
		require.Equal(t, d, cd)

		_ = s.Clean()
	})

	t.Run("put file larger than cache size", func(t *testing.T) {
		size := int64(1000)
		s, _ := store.New("/tmp/")
		c, _ := New(size, s)

		d, _ := os.ReadFile(testFiles[0].uri)
		err := c.Put(testFiles[0].uri, d)

		require.ErrorIs(t, err, ErrFileToLarge)

		_ = s.Clean()
	})

	t.Run("cache oversize", func(t *testing.T) {
		size := int64(300000)
		s, _ := store.New("/tmp/")
		c, _ := New(size, s)

		for _, file := range testFiles {
			d, _ := os.ReadFile(file.uri)

			err := c.Put(file.uri, d)

			require.Nil(t, err)
		}

		ds := getDirSize(s.Path())

		require.LessOrEqual(t, ds, c.size)

		_ = s.Clean()
	})
}

func getDirSize(path string) int64 {
	dir, _ := os.Open(path)
	defer dir.Close()
	files, _ := dir.Readdir(0)
	size := int64(0)
	for _, file := range files {
		size += file.Size()
	}

	return size
}
