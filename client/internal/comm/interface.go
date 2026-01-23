package comm

import (
	"github.com/boginskiy/GophKeeper/client/cmd/client"
	"github.com/boginskiy/GophKeeper/client/internal/user"
)

type Rooter interface {
	Execute(bool, *client.ClientCLI, *user.UserCLI)
}

type Commander interface {
	Registration(*client.ClientCLI, *user.UserCLI)
}
