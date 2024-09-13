package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/yakuninmax/imgpreviewer/internal/logger"
)

const (
	cacheSizeEnvName      = "IMPR_CACHE_SIZE"
	cacheFolderEnvName    = "IMPR_CACHE_FOLDER"
	requestTimeoutEnvName = "IMPR_REQ_TIMEOUT"
	defaultCacheSize      = 10485760
	defaultCachePath      = "/tmp/impr_cache"
	defaultRequestTimeout = 10
)

var (
	ErrCacheSizeZeroOrLess      = errors.New("cache size is zero or less")
	ErrRequestTimeoutZeroOrLess = errors.New("request timeout is zero or less")
)

type Config struct {
	cacheSize      int64
	cachePath      string
	requestTimeout time.Duration
}

func New(logger *logger.Log) (*Config, error) {
	cacheSize, err := getCacheSize(logger)
	if err != nil {
		return nil, err
	}

	cachePath, err := getCachePath(logger)
	if err != nil {
		return nil, err
	}

	requestTimeout, err := getRequestTimeout(logger)
	if err != nil {
		return nil, err
	}

	return &Config{
		cacheSize:      cacheSize,
		cachePath:      cachePath,
		requestTimeout: requestTimeout,
	}, nil
}

func (c *Config) CacheSize() int64 {
	return c.cacheSize
}

func (c *Config) CachePath() string {
	return c.cachePath
}

func (c *Config) RequestTimeout() time.Duration {
	return c.requestTimeout
}

// Get cache size from env var.
func getCacheSize(logger *logger.Log) (int64, error) {
	env := os.Getenv(cacheSizeEnvName)

	// Check if no env, or empty string.
	if env == "" {
		logger.Warn("IMPR_CACHE_SIZE value is empty, set default cache size " + strconv.Itoa(defaultCacheSize/1024/1024) + "MB")

		return defaultCacheSize, nil
	}

	// Convert string parameter.
	size, err := strconv.Atoi(env)
	if err != nil {
		return 0, fmt.Errorf("failed to set cache size: %w", err)
	}

	// Check if the cache size less or equal 0.
	if size <= 0 {
		return 0, ErrCacheSizeZeroOrLess
	}

	// Convert MB to bytes.
	logger.Info("cache size is " + env + "MB")

	// Convert megabytes to bytes, and return.
	return int64(size * 1024 * 1024), nil
}

// Get cache folder path from env var.
func getCachePath(logger *logger.Log) (string, error) {
	path := os.Getenv(cacheFolderEnvName)

	if path == "" {
		logger.Warn("IMPR_CACHE_PATH value is empty, set default cache path " + defaultCachePath)

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

// Get request timeout.
func getRequestTimeout(logger *logger.Log) (time.Duration, error) {
	env := os.Getenv(requestTimeoutEnvName)

	// Check if no env, or empty string.
	if env == "" {
		logger.Warn("IMPR_REQ_TIMEOUT value is empty, set default request timeout " + strconv.Itoa(defaultRequestTimeout) + " seconds")

		return time.Duration(defaultRequestTimeout), nil
	}

	// Convert string parameter.
	timeout, err := strconv.Atoi(env)
	if err != nil {
		return 0, fmt.Errorf("failed to set request timeout: %w", err)
	}

	// Check if the timeout less or equal 0.
	if timeout <= 0 {
		return 0, ErrRequestTimeoutZeroOrLess
	}

	// Convert MB to bytes.

	logger.Info("request timeout is " + env + " seconds")

	// Convert int to time.Duration, and return.
	return time.Duration(timeout), nil
}
