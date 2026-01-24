package service

import (
	"fmt"

	"github.com/boginskiy/GophKeeper/client/cmd/client"
	"github.com/boginskiy/GophKeeper/client/cmd/config"
	"github.com/boginskiy/GophKeeper/client/internal/api"
	"github.com/boginskiy/GophKeeper/client/internal/cli"
	"github.com/boginskiy/GophKeeper/client/internal/logg"
	"github.com/boginskiy/GophKeeper/client/internal/model"
	"github.com/boginskiy/GophKeeper/client/internal/user"
)

type TexterService struct {
	Cfg        config.Config
	Logg       logg.Logger
	Dialoger   cli.Dialoger
	ServiceAPI api.ServiceAPI
}

func NewTexterService(
	cfg config.Config,
	logger logg.Logger,
	dialog cli.Dialoger,
	serviceAPI api.ServiceAPI,
) *TexterService {

	return &TexterService{
		Cfg:        cfg,
		Logg:       logger,
		Dialoger:   dialog,
		ServiceAPI: serviceAPI,
	}
}

func (t *TexterService) Create(client *client.ClientCLI, user *user.UserCLI) {
	// Get info about type text.
	name, _ := t.Dialoger.GetSomeThing(client, user,
		fmt.Sprintf("%s\n\r%s",
			"What type of data do you want to create: \n\r\t credentials \n\r\t text \n\r\t card",
			"come back: back, need help: help, pass: enter"))

	// Enter text for saving.
	tx, _ := t.Dialoger.GetSomeThing(client, user, "Enter the text...")
	t.ServiceAPI.CreateText(user, *model.NewText(name, tx, user.User.Email))
}

func (t *TexterService) Read(client *client.ClientCLI, user *user.UserCLI) {
	name, _ := t.Dialoger.GetSomeThing(client, user,
		fmt.Sprintf("%s\n\r%s",
			"What type of data do you want to read: \n\r\t credentials \n\r\t text \n\r\t card",
			"come back: back, need help: help, pass: enter"))

	t.ServiceAPI.ReadText(user, model.Text{Name: name})
}

func (t *TexterService) Update(client *client.ClientCLI, user *user.UserCLI) {
	name, _ := t.Dialoger.GetSomeThing(client, user,
		fmt.Sprintf("%s\n\r%s",
			"What type of data do you want to update: \n\r\t credentials \n\r\t text \n\r\t card",
			"come back: back, need help: help, pass: enter"))

	tx, _ := t.Dialoger.GetSomeThing(client, user, "Enter the new text...")
	t.ServiceAPI.UpdateText(user, *model.NewText(name, tx, user.User.Email))
}

func (t *TexterService) Delete(client *client.ClientCLI, user *user.UserCLI) {

}

// типы хранимой инфы
// пары логин/пароль;
// произвольные текстовые данные;
// произвольные бинарные данные;
// данные банковских карт.
