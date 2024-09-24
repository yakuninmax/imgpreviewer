package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

const (
	cacheSizeEnv          = "IMPR_CACHE_SIZE"
	cachePathEnv          = "IMPR_CACHE_PATH"
	requestTimeoutEnv     = "IMPR_REQ_TIMEOUT"
	serverPort            = "IMPR_PORT"
	defaultSereverPort    = "8080"
	defaultCacheSize      = 10485760
	defaultCachePath      = "/tmp/impr_cache"
	defaultRequestTimeout = 10
)

var (
	ErrCacheSizeZeroOrLess      = errors.New("cache size is zero or less")
	ErrRequestTimeoutZeroOrLess = errors.New("request timeout is zero or less")
	ErrInvalidPort              = errors.New("invalid port number")
)

type logger interface {
	Info(string)
	Warn(string)
	Error(string)
	Debug(string)
}

type Config struct {
	cacheSize      int64
	cachePath      string
	requestTimeout time.Duration
	serverPort     string
}

func New(logg logger) (*Config, error) {
	cs, err := getCacheSize(logg)
	if err != nil {
		return nil, err
	}

	cp, err := getCachePath(logg)
	if err != nil {
		return nil, err
	}

	rt, err := getRequestTimeout(logg)
	if err != nil {
		return nil, err
	}

	sp, err := getServerPort(logg)
	if err != nil {
		return nil, err
	}

	return &Config{
		cacheSize:      cs,
		cachePath:      cp,
		requestTimeout: rt,
		serverPort:     sp,
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

func (c *Config) Port() string {
	return c.serverPort
}

// Get cache size from env var.
func getCacheSize(logg logger) (int64, error) {
	env := os.Getenv(cacheSizeEnv)

	// Check if no env, or empty string.
	if env == "" {
		logg.Warn("IMPR_CACHE_SIZE value is empty, set default cache size " +
			strconv.Itoa(defaultCacheSize/1024/1024) + "MB")

		return int64(defaultCacheSize), nil
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
	logg.Info("cache size is " + env + "MB")

	// Convert megabytes to bytes, and return.
	return int64(size * 1024 * 1024), nil
}

// Get cache folder path from env var.
func getCachePath(logg logger) (string, error) {
	path := os.Getenv(cachePathEnv)

	if path == "" {
		logg.Warn("IMPR_CACHE_PATH value is empty, set default cache path " + defaultCachePath)

		path = defaultCachePath
	}

	// Check if path is not absolute.
	if !filepath.IsAbs(path) {
		absp, err := filepath.Abs(path)
		print(absp)
		if err != nil {
			return "", fmt.Errorf("invalid path: %w", err)
		}

		path = absp
	}

	return path, nil
}

// Get request timeout.
func getRequestTimeout(logg logger) (time.Duration, error) {
	env := os.Getenv(requestTimeoutEnv)

	// Check if no env, or empty string.
	if env == "" {
		logg.Warn("IMPR_REQ_TIMEOUT value is empty, set default request timeout " +
			strconv.Itoa(defaultRequestTimeout) + " seconds")

		return time.Duration(defaultRequestTimeout) * time.Second, nil
	}

	// Convert string parameter.
	to, err := strconv.Atoi(env)
	if err != nil {
		return 0, fmt.Errorf("failed to set request timeout: %w", err)
	}

	// Check if the timeout less or equal 0.
	if to <= 0 {
		return 0, ErrRequestTimeoutZeroOrLess
	}

	// Convert MB to bytes.

	logg.Info("request timeout is " + env + " seconds")

	// Convert int to time.Duration, and return.
	return time.Duration(to) * time.Second, nil
}

// Get server port.
func getServerPort(logg logger) (string, error) {
	env := os.Getenv(serverPort)

	// Check if no env, or empty string.
	if env == "" {
		logg.Warn("IMPR_PORT value is empty, set default port " + defaultSereverPort)

		return defaultSereverPort, nil
	}

	// Check port number.
	p, err := strconv.Atoi(env)
	if err != nil {
		return "", err
	}

	if p < 1 || p > 65535 {
		return "", ErrInvalidPort
	}

	return env, nil
}
