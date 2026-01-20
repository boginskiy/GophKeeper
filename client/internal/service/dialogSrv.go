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
	// Hello
	d.Hello(client, user)

	// Identification, Authorization, Registration user.
	d.Auth.Identification(client, user)

	// Бесконечный диалог
	for {
		mess, err := user.ReceiveMess()
		d.Logger.CheckWithFatal(err, "reading error of cli")

		if mess == "exit" {
			break
		}

	}

	// Active action.
	fmt.Println(">> End CLI Session <<")

	// Save data current user in config.file for feature.
	// defer d.Auth.Identity.SaveCurrentUser(user)
}

func (d *DialogService) Hello(client *client.ClientCLI, user *user.UserCLI) {
	client.SendMess("Hello! Press the 'Enter'...")
	user.ReceiveMess()
}
