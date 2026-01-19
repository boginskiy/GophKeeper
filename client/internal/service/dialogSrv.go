package service

import (
	"fmt"

	"github.com/boginskiy/GophKeeper/client/cmd/client"
	"github.com/boginskiy/GophKeeper/client/internal/auth"
	"github.com/boginskiy/GophKeeper/client/internal/logg"
	"github.com/boginskiy/GophKeeper/client/internal/user"
)

type DialogService struct {
	Logger logg.Logger
	Auth   *auth.Auth
	Client *client.ClientCLI
	User   *user.UserCLI
}

func NewDialogService(logger logg.Logger, clt *client.ClientCLI, usr *user.UserCLI, auth *auth.Auth) *DialogService {
	return &DialogService{
		Logger: logger,
		Client: clt,
		User:   usr,
		Auth:   auth,
	}
}

func (d *DialogService) Run(client *client.ClientCLI, user *user.UserCLI) {
	// Identification, Authorization, Registration user.
	d.Auth.Identification(client, user)

	// Active action.
	fmt.Println(">> End CLI Session <<")

	// Save data current user in config.file for feature.
	defer d.Auth.Identity.SaveCurrentUser(user)
}
