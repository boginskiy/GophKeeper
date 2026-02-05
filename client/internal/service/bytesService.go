package service

import (
	"github.com/boginskiy/GophKeeper/client/cmd/config"
	"github.com/boginskiy/GophKeeper/client/internal/api"
	"github.com/boginskiy/GophKeeper/client/internal/infra"
	"github.com/boginskiy/GophKeeper/client/internal/logg"
	"github.com/boginskiy/GophKeeper/client/internal/model"
	"github.com/boginskiy/GophKeeper/client/internal/user"
	"github.com/boginskiy/GophKeeper/client/internal/utils"
)

type BytesService struct {
	Cfg         config.Config
	Logg        logg.Logger
	FileHandler utils.FileHandler
	PathHandler utils.PathHandler
	RemoteByter api.RemoteByter[model.Bytes]
	FileManager infra.FileManager
}

func NewBytesService(
	cfg config.Config,
	logger logg.Logger,
	fileHandler utils.FileHandler,
	pathHandler utils.PathHandler,
	remoteByter api.RemoteByter[model.Bytes],
	fileManager infra.FileManager,
) *BytesService {

	return &BytesService{
		Cfg:         cfg,
		Logg:        logger,
		FileHandler: fileHandler,
		PathHandler: pathHandler,
		RemoteByter: remoteByter,
		FileManager: fileManager,
	}
}

func (t *BytesService) Upload(user user.User, pathToFile string, tp string) (any, error) {
	bytes, err := t.FileManager.GetModelBytesFromFile(pathToFile, tp)
	if err != nil {
		return nil, err
	}
	defer bytes.Descr.Close()
	return t.RemoteByter.Upload(user, *bytes)
}

func (t *BytesService) Unload(user user.User, fileName string) (any, error) {
	modBytes := &model.Bytes{Name: fileName, Type: t.FileHandler.GetTypeFile(fileName)}

	file, path, err := t.FileHandler.CreateFileInStore(modBytes)
	if err != nil {
		return nil, err
	}

	defer file.Close()
	modBytes.Descr = file

	serverHeader, err := t.RemoteByter.Unload(user, *modBytes)

	if err != nil {
		t.FileHandler.DeleteFile(path)
		return nil, err
	}
	// В ответ кидаем заголовок. В нем вся инфа.
	return serverHeader, nil
}

func (t *BytesService) Read(user user.User, nameFile string) (any, error) {
	return t.RemoteByter.Read(user, model.Bytes{Name: nameFile})
}

func (t *BytesService) ReadAll(user user.User, typeFile string) (any, error) {
	return t.RemoteByter.ReadAll(user, model.Bytes{Type: typeFile})
}

func (t *BytesService) Delete(user user.User, nameFile string) (any, error) {
	return t.RemoteByter.Delete(user, model.Bytes{Name: nameFile})
}
