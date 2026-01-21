package service

import (
	"fmt"

	"github.com/boginskiy/GophKeeper/client/cmd/client"
	"github.com/boginskiy/GophKeeper/client/cmd/config"
	"github.com/boginskiy/GophKeeper/client/internal/auth"
	"github.com/boginskiy/GophKeeper/client/internal/cli"
	"github.com/boginskiy/GophKeeper/client/internal/logg"
	"github.com/boginskiy/GophKeeper/client/internal/user"
)

type KeeperService struct {
	Cfg        config.Config
	Logg       logg.Logger
	Identifier auth.Identifier
	Dialoger   cli.Dialoger
	Auth       *auth.Auth
}

func NewKeeperService(
	cfg config.Config,
	logger logg.Logger,
	identity auth.Identifier,
	dialoger cli.Dialoger,
	auth *auth.Auth) *KeeperService {

	return &KeeperService{
		Cfg:        cfg,
		Logg:       logger,
		Identifier: identity,
		Dialoger:   dialoger,
		Auth:       auth,
	}
}

func (d *KeeperService) Run(client *client.ClientCLI, user *user.UserCLI) {
	d.ExecuteHello(client, user) // Hello
	d.ExecuteAuth(client, user)  // Auth

	// Бесконечный диалог
	for {
		mess, err := user.ReceiveMess()
		d.Logg.CheckWithFatal(err, "reading error of cli")

		if mess == "exit" {
			break
		}

	}

	// Active action.
	fmt.Println(">> End CLI Session <<")

	// Save data current user in config.file for feature.
	// defer d.Auth.Identity.SaveCurrentUser(user)
}

func (d *KeeperService) ExecuteHello(client *client.ClientCLI, user *user.UserCLI) {
	d.Dialoger.ShowHello(client, user)
}

func (d *KeeperService) ExecuteAuth(client *client.ClientCLI, user *user.UserCLI) {
	// Identification user. load data from config.file if exist.
	IsThereRegistr := d.Identifier.Identification(user)

	// Authentication user.
	IsThereAuthent := d.Auth.Authentication(IsThereRegistr, client, user)

	// Registration user.
	d.Auth.Registration(IsThereAuthent, client, user)

}
