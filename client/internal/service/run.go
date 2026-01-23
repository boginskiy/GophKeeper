package service

import (
	"github.com/boginskiy/GophKeeper/client/cmd/client"
	"github.com/boginskiy/GophKeeper/client/cmd/config"
	"github.com/boginskiy/GophKeeper/client/internal/auth"
	"github.com/boginskiy/GophKeeper/client/internal/cli"
	"github.com/boginskiy/GophKeeper/client/internal/comm"
	"github.com/boginskiy/GophKeeper/client/internal/logg"
	"github.com/boginskiy/GophKeeper/client/internal/user"
)

type Run struct {
	Cfg        config.Config
	Logg       logg.Logger
	Identifier auth.Identifier
	Dialoger   cli.Dialoger
	Auth       *auth.Auth
	Root       comm.Rooter
}

func NewRun(
	cfg config.Config,
	logger logg.Logger,
	identity auth.Identifier,
	dialoger cli.Dialoger,
	auth *auth.Auth,
	root comm.Rooter) *Run {

	return &Run{
		Cfg:        cfg,
		Logg:       logger,
		Identifier: identity,
		Dialoger:   dialoger,
		Auth:       auth,
		Root:       root,
	}
}

func (d *Run) Run(client *client.ClientCLI, user *user.UserCLI) {
	d.Dialoger.ShowHello(client, user) // Hello
	ok := d.ExecuteAuth(client, user)  // Auth
	d.Root.Execute(ok, client, user)   // Root

	d.Dialoger.ShowSomeInfo(client,
		"Session is over. Goodbye")
	// Save data current user in config.file for feature.
	defer d.Identifier.SaveCurrentUser(user)
}

func (d *Run) ExecuteAuth(client *client.ClientCLI, user *user.UserCLI) bool {
	// Identification user. load data from config.file if exist.
	IsThereRegistr := d.Identifier.Identification(user)
	// Authentication user.
	IsThereAuthent := d.Auth.Authentication(IsThereRegistr, client, user)
	// Registration user.
	IsThereRegistr = d.Auth.Registration(IsThereAuthent, client, user)

	return IsThereRegistr || IsThereAuthent
}
