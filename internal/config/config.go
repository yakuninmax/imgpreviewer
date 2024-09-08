package config

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"

	"github.com/pkg/errors"
)

const (
	defaultCacheSize = 10485760
	defaultCachePath = "cache"
)

var (
	cacheSizeEnvName       = "IMPR_CACHE_SIZE"
	cacheFolderEnvName     = "IMPR_CACHE_FOLDER"
	ErrCacheSizeZeroOrLess = errors.New("cache size is zero or less")
)

type Config struct {
	cacheSize int64
	cachePath string
}

func New(logger *slog.Logger) (*Config, error) {
	cacheSize, err := getCacheSize(logger)
	if err != nil {
		return nil, err
	}

	cachePath, err := getCachePath(logger)
	if err != nil {
		return nil, err
	}

	return &Config{
		cacheSize: cacheSize,
		cachePath: cachePath,
	}, nil
}

func (c *Config) CacheSize() int64 {
	return c.cacheSize
}

func (c *Config) CachePath() string {
	return c.cachePath
}

// Get cache size from env var.
func getCacheSize(logger *slog.Logger) (int64, error) {
	env := os.Getenv(cacheSizeEnvName)

	// Check if no env, or empty string.
	if env == "" {
		logger.Warn("IMPR_CACHE_SIZE value is empty, set default cache size")
		logger.Info("cache size is " + strconv.Itoa(defaultCacheSize/1024/1024) + "MB")

		return defaultCacheSize, nil
	}

	// Convert string parameter.
	size, err := strconv.Atoi(env)
	if err != nil {
		return 0, fmt.Errorf("failed to set cache size: %w", err)
	}

	// Check if the cache size less than 1.
	if size <= 0 {
		return 0, ErrCacheSizeZeroOrLess
	}

	// Convert MB to bytes.

	logger.Info("cache size is " + env + "MB")

	// Convert megabytes to bytes, and return.
	return int64(size * 1024 * 1024), nil
}

// Get cache folder path from env var.
func getCachePath(logger *slog.Logger) (string, error) {
	path := os.Getenv(cacheFolderEnvName)

	if path == "" {
		logger.Warn("IMPR_CACHE_PATH value is empty, set default cache path")

		path = defaultCachePath
	}

	// Check if path is not absolute.
	if !filepath.IsAbs(path) {
		absPath, err := filepath.Abs(path)
		if err != nil {
			return "", fmt.Errorf("invalid path: %w", err)
		}

		path = absPath
	}

	return path, nil
}
