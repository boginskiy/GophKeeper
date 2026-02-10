package infra

import (
	"os"
	"strconv"

	"github.com/boginskiy/GophKeeper/client/cmd/config"
	"github.com/boginskiy/GophKeeper/client/internal/errs"
	"github.com/boginskiy/GophKeeper/client/internal/model"
	"github.com/boginskiy/GophKeeper/client/internal/utils"
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

func (f *FileService) GetModelBytesFromFile(pathToFile string, typ string) (*model.Bytes, error) {
	pathToFile, err := f.PathHandler.TransPathToAbs(pathToFile)

	if err != nil || !f.FileHandler.CheckOfFile(pathToFile) {
		return nil, errs.ErrReadFile.Wrap(err)
	}

	name := f.FileHandler.GetFileFromPath(pathToFile)
	descr, err := f.FileHandler.GetDescrFile(pathToFile)
	if err != nil {
		return nil, err
	}

	size, err := f.FileHandler.GetSizeFile(descr)
	if err != nil {
		return nil, err
	}

	return &model.Bytes{
		Name:     name,
		Descr:    descr,
		SentSize: strconv.FormatInt(size, 10),
		Type:     typ,
	}, nil
}

func (f *FileService) CreateFileInStore(modBytes *model.Bytes) (file *os.File, path string, err error) {
	// Create path
	path, err = f.PathHandler.CreatePathToWd(config.STORE, modBytes.Type, modBytes.Name)
	if err != nil {
		return nil, "", err
	}

	// Create all folders
	err = f.FileHandler.CreateFolder(path, MOD)
	if err != nil {
		return nil, "", err
	}

	// Create file
	file, err = f.FileHandler.ReadOrCreateFile(path, MOD)
	if err != nil {
		return nil, "", err
	}
	return file, path, err
}
