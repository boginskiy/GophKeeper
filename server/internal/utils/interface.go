package utils

import "os"

type PathCreater interface {
	CreatePathToConfig(nameApp, nameFile string) (path string, err error)
	CreatePathToWd(f1, f2, n string) (path string, err error)
}

type Creater interface {
	CreateFolder(path string, mod os.FileMode) error
	ReadOrCreateFile(path string, mod os.FileMode) (file *os.File, err error)
}

type Serializater interface {
	Deserialization(data []byte, obj any) error
	Serialization(obj any) ([]byte, error)
}

type Cleaner interface {
	TruncateFile(path string, mod os.FileMode) (file *os.File, err error)
}

type Pather interface {
	TakeDescrFromFile(string) (*os.File, error)
	TransPathToAbs(string) (string, error)
	TakeFileFromPath(string) string
}

type FileChecker interface {
	CheckOfFile(string) bool
}

type FileHandler interface {
	Serializater
	FileChecker
	PathCreater
	Creater
	Cleaner
	Pather

	TakeSizeFile(*os.File) (int64, error)
	GetTypeFile(string) string
}

type Checker interface {
	CheckTwoString(oneStr, twoStr string) bool
	CheckPassword(userPassword, password string) bool
}
