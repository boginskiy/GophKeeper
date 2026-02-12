package utils

import "os"

// FileHandler
type FileHandler interface {
	FileChecker
	FileCreater
	FileCleaner
	FileGetter
}

// Checker
type Checker interface {
	CheckTwoString(oneStr, twoStr string) bool
	CheckPassword(userPassword, password string) bool
}

// PathHandler
type PathHandler interface {
	PathCreater
	TransPathToAbs(string) (string, error)
}

type PathCreater interface {
	CreatePathToConfig(elem ...string) (path string, err error)
	CreatePathToWd(elem ...string) (path string, err error)
}

type FileCreater interface {
	CreateFolder(path string, mod os.FileMode) error
	ReadOrCreateFile(path string, mod os.FileMode) (file *os.File, err error)
}

type FileCleaner interface {
	TruncateFile(path string, mod os.FileMode) (file *os.File, err error)
}

type FileGetter interface {
	GetDescrFile(string) (*os.File, error)
	GetFileFromPath(string) string
	GetSizeFile(*os.File) (int64, error)
	GetTypeFile(string) string
}

type FileChecker interface {
	CheckOfFile(string) bool
}

type Serializater interface {
	Deserialization(data []byte, obj any) error
	Serialization(obj any) ([]byte, error)
}
