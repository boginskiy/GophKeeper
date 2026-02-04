package comm

import (
	"github.com/boginskiy/GophKeeper/client/cmd/client"
	"github.com/boginskiy/GophKeeper/client/internal/auth"
	"github.com/boginskiy/GophKeeper/client/internal/user"
)

type Rooter interface {
	ExecuteComm(bool, *client.ClientCLI, *user.UserCLI)
	ExecuteAuth(auth.Auth, user.User) bool
}

type Commander interface {
	Registration(user.User, string)
}
