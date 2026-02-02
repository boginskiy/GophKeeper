package service

import (
	"os"

	"github.com/boginskiy/GophKeeper/client/cmd/config"
	"github.com/boginskiy/GophKeeper/client/internal/api"
	"github.com/boginskiy/GophKeeper/client/internal/logg"
	"github.com/boginskiy/GophKeeper/client/internal/model"
	"github.com/boginskiy/GophKeeper/client/internal/user"
	"github.com/boginskiy/GophKeeper/client/internal/utils"
)

type ByterService struct {
	Cfg         config.Config
	Logg        logg.Logger
	FileHandler utils.FileHandler
	ServiceAPI  api.ServiceAPI
}

func NewByterService(
	cfg config.Config,
	logger logg.Logger,
	fileHdlr utils.FileHandler,
	serviceAPI api.ServiceAPI,
) *ByterService {

	return &ByterService{
		Cfg:         cfg,
		Logg:        logger,
		FileHandler: fileHdlr,
		ServiceAPI:  serviceAPI,
	}
}

func (t *ByterService) Upload(user user.User, pathToFile string) (any, error) {
	bytes, err := model.NewBytesFromFile(t.FileHandler, pathToFile)
	if err != nil {
		return nil, err
	}
	defer bytes.Descr.Close()

	return t.ServiceAPI.Upload(user, *bytes)
}

func (t *ByterService) Unload(user user.User, fileName string) (any, error) {
	modBytes := &model.Bytes{Name: fileName, Type: t.FileHandler.GetTypeFile(fileName)}

	file, path, err := t.FileHandler.CreateFileInStore(modBytes)
	if err != nil {
		return nil, err
	}

	defer file.Close()
	modBytes.Descr = file

	_, err = t.ServiceAPI.Unload(user, *modBytes)
	if err != nil {
		// Удаляем созданный ранее файл
		os.Remove(path)
		return nil, err
	}

	// Что передать в хендлер ?
	return nil, nil
}

// TODO!
// Шифрование: Если безопасность критична, подумайте о шифровании файлов до отправки и дешифрации на стороне сервера.
// Контроль целостности: Возможно добавить проверку контрольных сумм (CRC, SHA-256) для гарантированной доставки файлов без повреждений.
