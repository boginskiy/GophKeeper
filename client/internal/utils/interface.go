package utils

import "os"

type Creater interface {
	CreateFolder(configPath string, mod os.FileMode) error
	CreatePathToConfig(nameApp, nameFile string) (path string, err error)
	ReadOrCreateFile(configPath string, mod os.FileMode) (file *os.File, err error)
}

type Serializater interface {
	Deserialization(data []byte, obj any) error
	Serialization(obj any) ([]byte, error)
}

type Cleaner interface {
	TruncateFile(configPath string, mod os.FileMode) (file *os.File, err error)
}

type Pather interface {
	TakeDescrFromFile(string) (*os.File, error)
	TransPathToAbs(string) (string, error)
	TakeFileFromPath(string) string
	DeleteFile(path string) error
}

type FileChecker interface {
	CheckOfFile(string) bool
}

type FileHandler interface {
	Serializater
	FileChecker
	Creater
	Cleaner
	Pather

	CreateFileInStore(obj PathCreater) (file *os.File, path string, err error)
	TakeSizeFile(*os.File) (int64, error)
	GetTypeFile(fileName string) string
}

type PathCreater interface {
	GetFileType() string
	GetFileName() string
}
