package service

import (
	"github.com/boginskiy/GophKeeper/client/cmd/config"
	"github.com/boginskiy/GophKeeper/client/internal/api"
	"github.com/boginskiy/GophKeeper/client/internal/logg"
	"github.com/boginskiy/GophKeeper/client/internal/model"
	"github.com/boginskiy/GophKeeper/client/internal/user"
	"github.com/boginskiy/GophKeeper/client/internal/utils"
)

type BytesService struct {
	Cfg         config.Config
	Logg        logg.Logger
	FileHandler utils.FileHandler
	RemoteByter api.RemoteByter[model.Bytes]
}

func NewBytesService(
	cfg config.Config,
	logger logg.Logger,
	fileHdlr utils.FileHandler,
	remoteByter api.RemoteByter[model.Bytes],
) *BytesService {

	return &BytesService{
		Cfg:         cfg,
		Logg:        logger,
		FileHandler: fileHdlr,
		RemoteByter: remoteByter,
	}
}

func (t *BytesService) Upload(user user.User, pathToFile string, tp string) (any, error) {
	bytes, err := model.NewBytesFromFile(t.FileHandler, pathToFile, tp)
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
