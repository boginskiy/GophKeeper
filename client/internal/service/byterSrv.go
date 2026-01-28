package service

import (
	"fmt"

	"github.com/boginskiy/GophKeeper/client/cmd/client"
	"github.com/boginskiy/GophKeeper/client/cmd/config"
	"github.com/boginskiy/GophKeeper/client/internal/api"
	"github.com/boginskiy/GophKeeper/client/internal/cli"
	"github.com/boginskiy/GophKeeper/client/internal/errs"
	"github.com/boginskiy/GophKeeper/client/internal/logg"
	"github.com/boginskiy/GophKeeper/client/internal/model"
	"github.com/boginskiy/GophKeeper/client/internal/rpc"
	"github.com/boginskiy/GophKeeper/client/internal/user"
	"github.com/boginskiy/GophKeeper/client/internal/utils"
)

type ByterService struct {
	Cfg         config.Config
	Logg        logg.Logger
	FileHandler utils.FileHandler
	Dialoger    cli.Dialoger
	ServiceAPI  api.ServiceAPI
	Type        string
}

func NewByterService(
	cfg config.Config,
	logger logg.Logger,
	fileHdlr utils.FileHandler,
	dialog cli.Dialoger,
	serviceAPI api.ServiceAPI,
) *ByterService {

	return &ByterService{
		Cfg:         cfg,
		Logg:        logger,
		FileHandler: fileHdlr,
		Dialoger:    dialog,
		ServiceAPI:  serviceAPI,
		Type:        "bytes",
	}
}

func (t *ByterService) Upload(client *client.ClientCLI, user *user.UserCLI) {
	// Enter text for saving.
	pathToFile, _ := t.Dialoger.GetSomeThing(client, user, "Enter the abs path to file...")

	bytes, err := model.NewBytesFromFile(t.FileHandler, pathToFile)
	if err != nil {
		t.Dialoger.ShowErr(client, err)
		return
	}

	defer bytes.Descr.Close()

	// API.
	obj, err := t.ServiceAPI.Upload(user, *bytes)
	if err != nil {
		t.Dialoger.ShowErr(client, err)
		return
	}

	// Conversion.
	res, ok := obj.(*rpc.UploadBytesResponse)
	if !ok {
		t.Dialoger.ShowErr(client, errs.ErrTypeConversion)
		return
	}

	// Response.
	t.Dialoger.ShowIt(client, fmt.Sprintf("%s %s\n\r", res.Status, res.UpdatedAt))
}

func (t *ByterService) Unload(client *client.ClientCLI, user *user.UserCLI) {

}

// TODO!
// Шифрование: Если безопасность критична, подумайте о шифровании файлов до отправки и дешифрации на стороне сервера.
// Контроль целостности: Возможно добавить проверку контрольных сумм (CRC, SHA-256) для гарантированной доставки файлов без повреждений.

// Читаем файл и получаем байты
// Наверно надо читать большие файлы и передавать их
