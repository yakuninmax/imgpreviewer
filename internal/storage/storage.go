package storage

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
)

var (
	ErrPathIsNotDir = errors.New("given path is not a dir")
)

type Storage struct {
	path  string
	limit int
}

func New(path string, cacheSize int, logger *slog.Logger) (*Storage, error) {
	tempDirPath, err := createFolder(path, logger)
	if err != nil {
		return nil, err
	}

	return &Storage{
		path:  tempDirPath,
		limit: cacheSize,
	}, nil
}

// Get storage path.
func (s *Storage) Path() string {
	return s.path
}

// Write file to storage.
func (s *Storage) Write(name string, data []byte) error {
	err := os.WriteFile(name, data, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// Read file from storage.
func (s *Storage) Read(name string) ([]byte, error) {
	data, err := os.ReadFile(name)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return data, nil
}

// Delete file from storage.
func (s *Storage) Delete(name string) error {
	return nil
}

// Remove storage temp dir.
func (s *Storage) Clean() error {
	err := os.RemoveAll(s.path)
	if err != nil {
		return fmt.Errorf("failed to remove cache dir: %w", err)
	}

	return nil
}

// Create cache temp folder.
func createFolder(path string, logger *slog.Logger) (string, error) {
	// Get temp dir name for cache at given path.
	// Use current date as temp dir name.
	date := time.Now().Format("20060102150405")
	tempDirPath := filepath.Join(path, date)

	// Check if given path exists.
	stat, err := os.Stat(path)

	// If path not exists, or exists and it is dir, create temp dir.
	if errors.Is(err, os.ErrNotExist) || stat.IsDir() {
		if err := os.MkdirAll(tempDirPath, os.ModePerm); err != nil {
			return "", fmt.Errorf("failed to create cache dir: %w", err)
		}
	} else {
		if !stat.IsDir() {
			return "", ErrPathIsNotDir
		}

		return "", fmt.Errorf("failed to get dir: %w", err)
	}

	logger.Info("cache path is " + tempDirPath)

	return tempDirPath, nil
}
