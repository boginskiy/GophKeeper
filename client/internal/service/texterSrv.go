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
	name := t.Dialoger.DialogsAbAction(client, user, "create")            // Get info about name text.
	text, _ := t.Dialoger.GetSomeThing(client, user, "Enter the text...") // Enter text for saving.

	// Call service.
	obj, err := t.ServiceAPI.Create(user, *model.NewText(name, t.Type, text, user.User.Email))

	// Proc errors.
	if err != nil {
		t.Dialoger.ShowErr(client, err)
		return
	}

	res, ok := obj.(*rpc.CreateResponse)
	if !ok {
		t.Dialoger.ShowErr(client, errs.ErrTypeConversion)
		return
	}

	// Response in cli.
	t.Dialoger.ShowIt(client, fmt.Sprintf("%s %s\n\r", res.Status, res.UpdatedAt))
}

func (t *TexterService) Read(client *client.ClientCLI, user *user.UserCLI) {
	name := t.Dialoger.DialogsAbAction(client, user, "read")

	obj, err := t.ServiceAPI.Read(user, model.Text{Name: name, Owner: user.User.Email})
	if err != nil {
		t.Dialoger.ShowErr(client, err)
		return
	}

	res, ok := obj.(*rpc.ReadResponse)
	if !ok {
		t.Dialoger.ShowErr(client, errs.ErrTypeConversion)
		return
	}

	t.Dialoger.ShowIt(client, fmt.Sprintf(
		"%s: %s   last update: %s\n\r",
		res.Name, res.Text, res.UpdatedAt))
}

func (t *TexterService) ReadAll(client *client.ClientCLI, user *user.UserCLI) {
	obj, err := t.ServiceAPI.ReadAll(user, model.Text{Type: t.Type, Owner: user.User.Email})
	if err != nil {
		t.Dialoger.ShowErr(client, err)
		return
	}

	res, ok := obj.(*rpc.ReadAllResponse)
	if !ok {
		t.Dialoger.ShowErr(client, errs.ErrTypeConversion)
		return
	}

	for _, text := range res.TextResponses {
		t.Dialoger.ShowIt(client, fmt.Sprintf(
			"%s: %s   last update: %s",
			text.Name, text.Text, text.UpdatedAt))
	}
	fmt.Println()
}

func (t *TexterService) Update(client *client.ClientCLI, user *user.UserCLI) {
	name := t.Dialoger.DialogsAbAction(client, user, "update")
	text, _ := t.Dialoger.GetSomeThing(client, user, "Enter the new text...")

	obj, err := t.ServiceAPI.Update(user, model.Text{Name: name, Tx: text, Owner: user.User.Email})

	if err != nil {
		t.Dialoger.ShowErr(client, err)
		return
	}

	res, ok := obj.(*rpc.CreateResponse)
	if !ok {
		t.Dialoger.ShowErr(client, errs.ErrTypeConversion)
		return
	}

	t.Dialoger.ShowIt(client, fmt.Sprintf("%s %s\n\r", res.Status, res.UpdatedAt))
}

func (t *TexterService) Delete(client *client.ClientCLI, user *user.UserCLI) {
	name := t.Dialoger.DialogsAbAction(client, user, "delete")

	obj, err := t.ServiceAPI.Delete(user, model.Text{Name: name, Owner: user.User.Email})

	if err != nil {
		t.Dialoger.ShowErr(client, err)
		return
	}

	res, ok := obj.(*rpc.DeleteResponse)
	if !ok {
		t.Dialoger.ShowErr(client, errs.ErrTypeConversion)
		return
	}
	t.Dialoger.ShowIt(client, fmt.Sprintf("%s %s\n\r", res.Status, res.Name))
}
