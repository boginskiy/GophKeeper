package infra

import (
	"os"

	"github.com/boginskiy/GophKeeper/server/cmd/config"
	"github.com/boginskiy/GophKeeper/server/internal/utils"
)

const MOD = 0755

type FileService struct {
	FileHandler utils.FileHandler
	PathHandler utils.PathHandler
}

func NewFileService(fileHandler utils.FileHandler, pathHandler utils.PathHandler) *FileService {
	return &FileService{
		FileHandler: fileHandler,
		PathHandler: pathHandler,
	}
}

func (m *FileService) CreateFileInStore(obj PathCreater) (file *os.File, path string, err error) {
	// Create path
	path, err = m.PathHandler.CreatePathToWd(config.STORE, obj.GetOwner(), obj.GetFileName())
	if err != nil {
		return nil, "", err
	}

	// Create all folders
	err = m.FileHandler.CreateFolder(path, MOD)
	if err != nil {
		return nil, "", err
	}

	// Create file
	file, err = m.FileHandler.ReadOrCreateFile(path, MOD)
	if err != nil {
		return nil, "", err
	}
	return file, path, err
}
