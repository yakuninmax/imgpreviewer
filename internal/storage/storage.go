package storage

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"strings"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz0123456789"

var ErrPathIsNotDir = errors.New("given path is not a dir")

type Storage struct {
	path string
}

func New(path string) (*Storage, error) {
	tmp, err := createFolder(path)
	if err != nil {
		return nil, err
	}

	return &Storage{
		path: tmp,
	}, nil
}

// Get storage path.
func (s *Storage) Path() string {
	return s.path
}

// Write file to storage.
func (s *Storage) Write(name string, data []byte) error {
	file := filepath.Join(s.path, name)

	err := os.WriteFile(file, data, os.ModePerm.Perm())
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// Read file from storage.
func (s *Storage) Read(name string) ([]byte, error) {
	file := filepath.Join(s.path, name)
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return data, nil
}

// Delete file from storage.
func (s *Storage) Delete(name string) error {
	filePath := filepath.Join(s.path, name)
	err := os.Remove(filePath)
	if err != nil {
		return err
	}

	return nil
}

// Remove storage temp dir.
func (s *Storage) Clean() error {
	err := os.RemoveAll(s.path)
	if err != nil {
		return fmt.Errorf("failed to remove temp cache folder: %w", err)
	}

	return nil
}

// Create cache temp folder.
func createFolder(path string) (string, error) {
	// Get dir name from time.
	dir, err := getRandomName()
	if err != nil {
		return "", err
	}

	tmp := filepath.Join(path, dir)

	// Check if given path exists.
	stat, err := os.Stat(path)

	// If path not exists, or exists and it is dir, create temp dir.
	if errors.Is(err, os.ErrNotExist) || stat.IsDir() {
		if err := os.MkdirAll(tmp, os.ModePerm); err != nil {
			return "", fmt.Errorf("failed to create cache dir: %w", err)
		}
	} else {
		if !stat.IsDir() {
			return "", ErrPathIsNotDir
		}

		return "", fmt.Errorf("failed to get dir: %w", err)
	}

	return tmp, nil
}

func getRandomName() (string, error) {
	var name strings.Builder

	for i := 0; i < 8; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(alphabet))))
		if err != nil {
			return "", err
		}

		name.WriteString(n.String())
	}

	return name.String(), nil
}
