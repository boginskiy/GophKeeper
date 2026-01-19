package auth

import (
	"github.com/boginskiy/GophKeeper/client/cmd/client"
	"github.com/boginskiy/GophKeeper/client/internal/model"
	"github.com/boginskiy/GophKeeper/client/internal/user"
)

type Auhter interface {
	Identification(client *client.ClientCLI, user *user.UserCLI)
	Registration(client *client.ClientCLI, user *user.UserCLI) bool
	Authentication(client *client.ClientCLI, user *user.UserCLI) bool
}

type Identification interface {
	Identification(*user.UserCLI) bool
	SaveCurrentUser(*model.User)
}
