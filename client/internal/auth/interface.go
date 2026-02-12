package auth

import (
	"github.com/boginskiy/GophKeeper/client/internal/model"
	"github.com/boginskiy/GophKeeper/client/internal/user"
)

type Auth interface {
	Authentication(user.User, *model.User) (string, error)
	Registration(user.User, *model.User) (string, error)
	Recovery(user.User, *model.User) (string, error)
	Identification(user.User) bool
}

type Identifier interface {
	Identification(user.User) bool
	SaveCurrentUser(user.User)
	Shutdown(<-chan struct{}, user.User)
}
