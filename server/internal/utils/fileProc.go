package utils

import (
	"os"
	"path/filepath"
	"strings"
)

type FileProc struct {
}

func NewFileProc() *FileProc {
	return &FileProc{}
}

func (f *FileProc) CreateFolder(path string, mod os.FileMode) error {
	// Creating a directory if it does not exist
	err := os.MkdirAll(filepath.Dir(path), mod)
	if err != nil {
		return err
	}
	return nil
}

func (f *FileProc) ReadOrCreateFile(path string, mod os.FileMode) (file *os.File, err error) {
	file, err = os.OpenFile(path, os.O_RDWR|os.O_CREATE, mod)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (f *FileProc) TruncateFile(path string, mod os.FileMode) (file *os.File, err error) {
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

func (f *FileProc) CheckOfFile(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func (f *FileProc) GetSizeFile(file *os.File) (int64, error) {
	stats, err := file.Stat()
	if err != nil {
		return 0, err
	}
	return stats.Size(), nil
}

func (f *FileProc) GetTypeFile(fileName string) string {
	res := strings.Split(fileName, ".")
	if len(res) > 1 {
		return res[len(res)-1]
	}
	return ""
}

// GetDescrFile
func (f *FileProc) GetDescrFile(path string) (*os.File, error) {
	return os.Open(path)
}

// GetFileFromPath
func (f *FileProc) GetFileFromPath(path string) string {
	arrPath := strings.Split(path, string(filepath.Separator))

	if len(arrPath) > 0 {
		return arrPath[len(arrPath)-1]
	}
	return ""
}
