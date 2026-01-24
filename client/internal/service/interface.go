package service

import (
	"github.com/boginskiy/GophKeeper/client/cmd/client"
	"github.com/boginskiy/GophKeeper/client/internal/user"
)

type Servicer interface {
	Create(*client.ClientCLI, *user.UserCLI)
	Read(*client.ClientCLI, *user.UserCLI)
	Update(*client.ClientCLI, *user.UserCLI)
	Delete(*client.ClientCLI, *user.UserCLI)
}
