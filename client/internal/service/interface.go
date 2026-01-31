package service

import (
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
	Upload(user.User, string) (any, error)
	Unload(user.User) (any, error)
}
