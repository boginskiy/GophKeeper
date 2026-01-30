package service

import (
	"github.com/boginskiy/GophKeeper/client/cmd/client"
	"github.com/boginskiy/GophKeeper/client/internal/auth"
	"github.com/boginskiy/GophKeeper/client/internal/user"
)

type Servicer[T any] interface { // ServicerTexter
	Create(user.User, T) (any, error)
	Read(user.User, T) (any, error)
	ReadAll(user.User, T) (any, error)
	Update(user.User, T) (any, error)
	Delete(user.User, T) (any, error)
}

type ServicerByter interface {
	Upload(*client.ClientCLI, *user.UserCLI)
	Unload(*client.ClientCLI, *user.UserCLI)
}

type Dialoger interface {
	DialogAbAuthentication(auth.Auth, user.User) bool
}
