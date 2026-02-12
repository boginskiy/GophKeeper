package utils

import (
	"os"
	"path/filepath"
)

type PathProc struct {
}

func NewPathProc() *PathProc {
	return &PathProc{}
}

// CreatePathToConfig create path to config file.
func (p *PathProc) CreatePathToConfig(elems ...string) (path string, err error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	elems = append([]string{dir}, elems...)
	return filepath.Join(elems...), nil
}

// CreatePathToWd
func (p *PathProc) CreatePathToWd(folders ...string) (path string, err error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	folders = append([]string{dir}, folders...)
	return filepath.Join(folders...), nil
}

// TransPathToAbs
func (p *PathProc) TransPathToAbs(path string) (string, error) {
	// Приведение пути к абсолютному виду.
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	// Нормализация.
	return filepath.Clean(absPath), nil
}
