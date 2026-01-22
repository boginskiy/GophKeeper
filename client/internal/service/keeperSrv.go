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
	messChan   chan string
	Identifier auth.Identifier
	Dialoger   cli.Dialoger
	Auth       *auth.Auth
}

func NewKeeperService(
	cfg config.Config,
	logger logg.Logger,
	ch chan string,
	identity auth.Identifier,
	dialoger cli.Dialoger,
	auth *auth.Auth) *KeeperService {

	return &KeeperService{
		Cfg:        cfg,
		Logg:       logger,
		messChan:   ch,
		Identifier: identity,
		Dialoger:   dialoger,
		Auth:       auth,
	}
}

func (d *KeeperService) Run(client *client.ClientCLI, user *user.UserCLI) {
	d.ExecuteHello(client, user)      // Hello
	ok := d.ExecuteAuth(client, user) // Auth

	// Бесконечный диалог
	for ok {
		comm, err := d.Dialoger.GetCommand(client, user)
		d.Logg.CheckWithFatal(err, "reading error of cli")

		if comm == "exit" {
			break

		} else if comm == "help" {
			d.ExecuteHelp(comm)

		} else {
			d.ExecuteCommand(comm)
		}
	}

	d.Dialoger.ShowSomeInfo(client, "Session is over. Goodbye")
	// Save data current user in config.file for feature.
	defer d.Identifier.SaveCurrentUser(user)
}

func (d *KeeperService) ExecuteHello(client *client.ClientCLI, user *user.UserCLI) {
	d.Dialoger.ShowHello(client, user)
}

func (d *KeeperService) ExecuteAuth(client *client.ClientCLI, user *user.UserCLI) bool {
	// Identification user. load data from config.file if exist.
	IsThereRegistr := d.Identifier.Identification(user)
	// Authentication user.
	IsThereAuthent := d.Auth.Authentication(IsThereRegistr, client, user)
	// Registration user.
	IsThereRegistr = d.Auth.Registration(IsThereAuthent, client, user)

	return IsThereRegistr || IsThereAuthent
}

func (d *KeeperService) ExecuteHelp(command string) {
	fmt.Println("Help")
}

func (d *KeeperService) ExecuteCommand(command string) {
	d.messChan <- command
	<-d.messChan
}
