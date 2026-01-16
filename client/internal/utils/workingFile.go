package utils

import (
	"encoding/json"
	"os"
	"path/filepath"
)

var NAMEFILE = "config.json"
var NAMEAPPL = "gophclient"

type WorkingFile struct {
}

func NewWorkingFile() *WorkingFile {
	return &WorkingFile{}
}

func (f *WorkingFile) CreateFolder(path string, mod os.FileMode) error {
	// Creating a directory if it does not exist
	err := os.MkdirAll(filepath.Dir(path), mod)
	if err != nil {
		return err
	}
	return nil
}

func (f *WorkingFile) CreatePathToConfig(nameApp, nameFile string) (path string, err error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, nameApp, nameFile), nil
}

func (f *WorkingFile) ReadOrCreateFile(path string, mod os.FileMode) (file *os.File, err error) {
	file, err = os.OpenFile(path, os.O_RDWR|os.O_CREATE, mod)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (f *WorkingFile) TruncateFile(path string, mod os.FileMode) (file *os.File, err error) {
	file, err = os.OpenFile(path, os.O_RDWR, mod)
	if err != nil {
		return nil, err
	}

	// Move it to the beginning of the file.
	_, err = file.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	// Clean old data.
	err = file.Truncate(0)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (f *WorkingFile) Deserialization(data []byte, obj any) error {
	return json.Unmarshal(data, obj)
}

func (f *WorkingFile) Serialization(obj any) ([]byte, error) {
	return json.MarshalIndent(obj, "", "  ")
}
