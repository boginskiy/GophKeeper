package app

import (
	"github.com/boginskiy/GophKeeper/client/cmd/client"
	"github.com/boginskiy/GophKeeper/client/cmd/config"
	"github.com/boginskiy/GophKeeper/client/internal/auth"
	"github.com/boginskiy/GophKeeper/client/internal/comm"
	"github.com/boginskiy/GophKeeper/client/internal/infra"
	"github.com/boginskiy/GophKeeper/client/internal/logg"
	"github.com/boginskiy/GophKeeper/client/internal/user"
)

type Runner struct {
	Cfg        config.Config
	Logg       logg.Logger
	Identifier auth.Identifier
	Dialoger   infra.Dialoger
	Root       comm.Rooter
	RootAuth   comm.Rooter
	done       chan struct{}
}

func NewRunner(
	cfg config.Config,
	logger logg.Logger,
	identity auth.Identifier,
	dialoger infra.Dialoger,
	root comm.Rooter,
	rootAuth comm.Rooter) *Runner {

	return &Runner{
		Cfg:        cfg,
		Logg:       logger,
		Identifier: identity,
		Dialoger:   dialoger,
		Root:       root,
		RootAuth:   rootAuth,
		done:       make(chan struct{}),
	}
}

func (d *Runner) Run(client *client.ClientCLI, user user.User) {
	d.Identifier.Shutdown(d.done, user)

	d.Dialoger.ShowIt("Hello, Man!")

	ok := d.RootAuth.ExecuteComm(true, user)
	d.Root.ExecuteComm(ok, user)

	d.Dialoger.ShowIt("Session is over. Goodbye")

	// Save data current user in config.file for feature.
	defer d.Identifier.SaveCurrentUser(user)
	defer close(d.done)
}
