package service

import (
	"github.com/boginskiy/GophKeeper/client/internal/user"
)

type TextServicer[T any] interface {
	Create(user.User, T) (any, error)
	Read(user.User, T) (any, error)
	ReadAll(user.User, T) (any, error)
	Update(user.User, T) (any, error)
	Delete(user.User, T) (any, error)
}

type BytesServicer interface {
	Upload(user.User, string) (any, error)
	Unload(user.User, string) (any, error)
	Read(user.User, string) (any, error)
	ReadAll(user.User, string) (any, error)
	Delete(user.User, string) (any, error)
}
