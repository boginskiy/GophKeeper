package infra

import "os"

type PathCreater interface {
	GetOwner() string
	GetFileName() string
}

type FileManager interface {
	CreateFileInStore(obj PathCreater) (file *os.File, path string, err error)
}
