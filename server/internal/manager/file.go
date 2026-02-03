package manager

import (
	"os"

	"github.com/boginskiy/GophKeeper/server/internal/utils"
)

const (
	STORE = "store"
	MOD   = 0755
)

type FileManage struct {
	FileHdler utils.FileHandler
}

func NewFileManage(fileHdler utils.FileHandler) *FileManage {
	return &FileManage{
		FileHdler: fileHdler,
	}
}

func (m *FileManage) CreateFileInStore(obj PathCreater) (file *os.File, path string, err error) {
	// Create path
	path, err = m.FileHdler.CreatePathToWd(STORE, obj.GetOwner(), obj.GetFileName())
	if err != nil {
		return nil, "", err
	}

	// Create all folders
	err = m.FileHdler.CreateFolder(path, MOD)
	if err != nil {
		return nil, "", err
	}

	// Create file
	file, err = m.FileHdler.ReadOrCreateFile(path, MOD)
	if err != nil {
		return nil, "", err
	}
	return file, path, err
}
