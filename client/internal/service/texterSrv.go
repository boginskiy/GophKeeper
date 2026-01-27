package service

import (
	"fmt"

	"github.com/boginskiy/GophKeeper/client/cmd/client"
	"github.com/boginskiy/GophKeeper/client/cmd/config"
	"github.com/boginskiy/GophKeeper/client/internal/api"
	"github.com/boginskiy/GophKeeper/client/internal/cli"
	"github.com/boginskiy/GophKeeper/client/internal/logg"
	"github.com/boginskiy/GophKeeper/client/internal/model"
	"github.com/boginskiy/GophKeeper/client/internal/rpc"
	"github.com/boginskiy/GophKeeper/client/internal/user"
)

type TexterService struct {
	Cfg        config.Config
	Logg       logg.Logger
	Dialoger   cli.Dialoger
	ServiceAPI api.ServiceAPI
	Type       string
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
		Type:       "text",
	}
}

func (t *TexterService) Create(client *client.ClientCLI, user *user.UserCLI) {
	// Get info about name text.
	name := t.Dialoger.DialogsAbAction(client, user, "create")

	// Enter text for saving.
	text, _ := t.Dialoger.GetSomeThing(client, user, "Enter the text...")

	// Call ServiceAPI.
	obj, err := t.ServiceAPI.Create(user, *model.NewText(name, t.Type, text, user.User.Email))

	if err != nil {
		t.Dialoger.ShowIt(client, err.Error())
	}

	res, ok := obj.(*rpc.CreateResponse)
	if !ok {
		t.Dialoger.ShowIt(client, "Type not valid")
	}

	t.Dialoger.ShowIt(client, fmt.Sprintf("%s: %s\n\r", res.Status, res.UpdatedAt))
	return

}

// func (t *TexterService) Read(client *client.ClientCLI, user *user.UserCLI) {
// 	name := t.Dialoger.DialogsAbAction(client, user, "read")

// 	t.ServiceAPI.Read(user, model.Text{Name: name, Type: t.Tp})
// }

// func (t *TexterService) Update(client *client.ClientCLI, user *user.UserCLI) {
// 	name := t.Dialoger.DialogsAbAction(client, user, "update")

// 	tx, _ := t.Dialoger.GetSomeThing(client, user, "Enter the new text...")
// 	t.ServiceAPI.Update(user, *model.NewText(name, t.Tp, tx, user.User.Email))
// }

// func (t *TexterService) Delete(client *client.ClientCLI, user *user.UserCLI) {

// }

// типы хранимой инфы
// пары логин/пароль;
// произвольные текстовые данные;
// произвольные бинарные данные;
// данные банковских карт.
