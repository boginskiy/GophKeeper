package cmd

import (
	"github.com/boginskiy/GophKeeper/client/cmd/client"
	"github.com/boginskiy/GophKeeper/client/cmd/config"
	"github.com/boginskiy/GophKeeper/client/internal/auth"
	"github.com/boginskiy/GophKeeper/client/internal/cli"
	"github.com/boginskiy/GophKeeper/client/internal/comm"
	"github.com/boginskiy/GophKeeper/client/internal/logg"
	"github.com/boginskiy/GophKeeper/client/internal/user"
)

type Runner struct {
	Cfg        config.Config
	Logg       logg.Logger
	Identifier auth.Identifier
	DialogSrv  cli.ShowGetter
	AuthSrv    auth.Auth
	Root       comm.Rooter
}

func NewRunner(
	cfg config.Config,
	logger logg.Logger,
	identity auth.Identifier,
	dialog cli.ShowGetter,
	authSrv auth.Auth,
	root comm.Rooter) *Runner {

	return &Runner{
		Cfg:        cfg,
		Logg:       logger,
		Identifier: identity,
		DialogSrv:  dialog,
		AuthSrv:    authSrv,
		Root:       root,
	}
}

func (d *Runner) Run(client *client.ClientCLI, user *user.UserCLI) {
	d.DialogSrv.ShowIt("Hello, Man!")

	ok := d.Root.ExecuteAuth(d.AuthSrv, user)
	d.Root.ExecuteComm(ok, client, user)

	d.DialogSrv.ShowIt("Session is over. Goodbye")

	// Save data current user in config.file for feature.
	defer d.Identifier.SaveCurrentUser(user)
}
