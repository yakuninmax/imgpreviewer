package storage

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	tfn := "testfile"

	s, err := New("/tmp")
	require.NoError(t, err)

	_, err = os.Stat(s.path)
	require.NotErrorIs(t, err, os.ErrNotExist)

	t.Run("write file", func(t *testing.T) {
		data, err := os.ReadFile("../../examples/_gopher_original_1024x504.jpg")
		require.NoError(t, err)

		err = s.Write(tfn, data)
		require.NoError(t, err)
	})

	t.Run("read file", func(t *testing.T) {
		orig, err := os.ReadFile("../../examples/_gopher_original_1024x504.jpg")
		require.NoError(t, err)

		data, err := s.Read(tfn)
		require.NoError(t, err)
		require.Equal(t, orig, data)
	})

	t.Run("delete file", func(t *testing.T) {
		err := s.Delete(tfn)
		require.NoError(t, err)
	})

	t.Run("clean storage", func(t *testing.T) {
		err := s.Clean()
		require.NoError(t, err)

		_, err = os.Stat(s.path)
		require.ErrorIs(t, err, os.ErrNotExist)
	})
}
