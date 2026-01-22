package auth

import (
	"github.com/boginskiy/GophKeeper/client/cmd/client"
	"github.com/boginskiy/GophKeeper/client/internal/user"
)

type Auhter interface {
	Registration(client *client.ClientCLI, user *user.UserCLI) bool
	Authentication(client *client.ClientCLI, user *user.UserCLI) bool
}

type Identifier interface {
	Identification(*user.UserCLI) bool
	SaveCurrentUser(*user.UserCLI)
}
