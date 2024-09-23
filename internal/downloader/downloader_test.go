package downloader

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	dl := New(10 * time.Second)
	var hdr map[string][]string

	t.Run("download image", func(t *testing.T) {
		orig, err := os.ReadFile("../../examples/space.jpg")
		require.NoError(t, err)

		img, err := dl.GetImage(
			"https://www.fileformat.info/format/jpeg/sample/0c047d42fdfb419e86c594f0f7ad3ce1/SPACE.JPG",
			hdr)
		require.NoError(t, err)
		require.Equal(t, orig, img)
	})

	t.Run("not image", func(t *testing.T) {
		_, err := dl.GetImage(
			"https://raw.githubusercontent.com/OtusGolang/final_project/refs/heads/master/03-image-previewer.md",
			hdr)
		require.ErrorIs(t, err, ErrInvalidFileType)
	})

	t.Run("remote server error", func(t *testing.T) {
		_, err := dl.GetImage(
			"https://raw.githubusercontent.com/OtusGolang/final_project/refs/heads/master/fake.file",
			hdr)
		require.EqualError(t, err, "remote server return: 404 Not Found")
	})
}
