package utils

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

const (
	NAMEFILE = "config.json"
	NAMEAPPL = "gophclient"
	STORE    = "store"
	MOD      = 0755
)

type FileHdlr struct {
}

func NewFileHdlr() *FileHdlr {
	return &FileHdlr{}
}

// CreateFolder
func (f *FileHdlr) CreateFolder(path string, mod os.FileMode) error {
	// Creating a directory if it does not exist
	err := os.MkdirAll(filepath.Dir(path), mod)
	if err != nil {
		return err
	}
	return nil
}

func (f *FileHdlr) CreatePathToConfig(nameApp, nameFile string) (path string, err error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, nameApp, nameFile), nil
}

// ReadOrCreateFile
func (f *FileHdlr) ReadOrCreateFile(path string, mod os.FileMode) (file *os.File, err error) {
	file, err = os.OpenFile(path, os.O_RDWR|os.O_CREATE, mod)
	if err != nil {
		return nil, err
	}
	return file, nil
}

// CreatePathToWd
func (f *FileHdlr) CreatePathToWd(f1, f2, f3 string) (path string, err error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, f1, f2, f3), nil
}

func (f *FileHdlr) TruncateFile(path string, mod os.FileMode) (file *os.File, err error) {
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

func (f *FileHdlr) Deserialization(data []byte, obj any) error {
	return json.Unmarshal(data, obj)
}

func (f *FileHdlr) Serialization(obj any) ([]byte, error) {
	return json.MarshalIndent(obj, "", "  ")
}

func (f *FileHdlr) TransPathToAbs(path string) (string, error) {
	// Приведение пути к абсолютному виду.
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	// Нормализация.
	return filepath.Clean(absPath), nil
}

func (f *FileHdlr) CheckOfFile(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func (f *FileHdlr) TakeFileFromPath(path string) string {
	arrPath := strings.Split(path, string(filepath.Separator))

	if len(arrPath) > 0 {
		return arrPath[len(arrPath)-1]
	}
	return ""
}

func (f *FileHdlr) TakeDescrFromFile(path string) (*os.File, error) {
	return os.Open(path)
}

func (f *FileHdlr) TakeSizeFile(file *os.File) (int64, error) {
	stats, err := file.Stat()
	if err != nil {
		return 0, err
	}
	return stats.Size(), nil
}

func (f *FileHdlr) GetTypeFile(fileName string) string {
	res := strings.Split(fileName, ".")
	if len(res) > 1 {
		return res[len(res)-1]
	}
	return ""
}

func (m *FileHdlr) CreateFileInStore(obj PathCreater) (file *os.File, path string, err error) {
	// Create path
	path, err = m.CreatePathToWd(STORE, obj.GetFileType(), obj.GetFileName())
	if err != nil {
		return nil, "", err
	}

	// Create all folders
	err = m.CreateFolder(path, MOD)
	if err != nil {
		return nil, "", err
	}

	// Create file
	file, err = m.ReadOrCreateFile(path, MOD)
	if err != nil {
		return nil, "", err
	}
	return file, path, err
}

func (m *FileHdlr) DeleteFile(path string) error {
	return os.Remove(path)
}
