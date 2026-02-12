package utils

import "os"

type FileHandler interface {
	FileChecker
	FileCreater
	FileCleaner
	FileDeleter
	FileGetter
}

// PathHandler
type PathHandler interface {
	PathCreater
	TransPathToAbs(string) (string, error)
}

// Serializater
type Serializater interface {
	Deserialization(data []byte, obj any) error
	Serialization(obj any) ([]byte, error)
}

type PathCreater interface {
	CreatePathToConfig(elem ...string) (path string, err error)
	CreatePathToWd(elem ...string) (path string, err error)
}

type FileDeleter interface {
	DeleteFile(path string) error
}

type FileCreater interface {
	CreateFolder(configPath string, mod os.FileMode) error
	ReadOrCreateFile(configPath string, mod os.FileMode) (file *os.File, err error)
}

type FileCleaner interface {
	TruncateFile(configPath string, mod os.FileMode) (file *os.File, err error)
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
