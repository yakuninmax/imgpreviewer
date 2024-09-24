package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	l "github.com/yakuninmax/imgpreviewer/internal/logger"
)

func TestConfig(t *testing.T) {
	logg, err := l.New()
	require.NoError(t, err)

	t.Run("get default values", func(t *testing.T) {
		os.Unsetenv("IMPR_CACHE_SIZE")
		os.Unsetenv("IMPR_CACHE_PATH")
		os.Unsetenv("IMPR_REQ_TIMEOUT")
		os.Unsetenv("IMPR_PORT")

		conf, err := New(logg)
		require.NoError(t, err)
		require.Equal(t, defaultCachePath, conf.cachePath)
		require.Equal(t, int64(defaultCacheSize), conf.cacheSize)
		require.Equal(t, defaultRequestTimeout*time.Second, conf.requestTimeout)
		require.Equal(t, defaultSereverPort, conf.serverPort)
	})

	t.Run("set values", func(t *testing.T) {
		os.Setenv("IMPR_CACHE_SIZE", "100")
		os.Setenv("IMPR_CACHE_PATH", "/tmp/test123")
		os.Setenv("IMPR_REQ_TIMEOUT", "60")
		os.Setenv("IMPR_PORT", "48080")

		conf, err := New(logg)
		require.NoError(t, err)
		require.Equal(t, "/tmp/test123", conf.cachePath)
		require.Equal(t, int64(100*1024*1024), conf.cacheSize)
		require.Equal(t, 60*time.Second, conf.requestTimeout)
		require.Equal(t, "48080", conf.serverPort)

		os.Unsetenv("IMPR_CACHE_SIZE")
		os.Unsetenv("IMPR_CACHE_PATH")
		os.Unsetenv("IMPR_REQ_TIMEOUT")
		os.Unsetenv("IMPR_PORT")
	})

	t.Run("invalid request timeout", func(t *testing.T) {
		os.Setenv("IMPR_CACHE_SIZE", "-100")

		_, err := New(logg)
		require.ErrorIs(t, ErrCacheSizeZeroOrLess, err)

		os.Unsetenv("IMPR_CACHE_SIZE")
	})

	t.Run("invalid request timeout", func(t *testing.T) {
		os.Setenv("IMPR_REQ_TIMEOUT", "-100")

		_, err := New(logg)
		require.ErrorIs(t, ErrRequestTimeoutZeroOrLess, err)

		os.Unsetenv("IMPR_REQ_TIMEOUT")
	})

	t.Run("invalid port", func(t *testing.T) {
		os.Setenv("IMPR_PORT", "76000")

		_, err := New(logg)
		require.ErrorIs(t, ErrInvalidPort, err)

		os.Unsetenv("IMPR_PORT")
	})
}
