package api

import (
	"github.com/boginskiy/GophKeeper/client/internal/model"
	"github.com/boginskiy/GophKeeper/client/internal/user"
)

type CRUD[T any] interface {
	Create(user.User, T) (any, error)
	Read(user.User, T) (any, error)
	Update(user.User, T) (any, error)
	Delete(user.User, T) (any, error)
}

type ReadAll[T any] interface {
	ReadAll(user.User, T) (any, error)
}

type Loader[T any] interface {
	Upload(user.User, T) (any, error)
	Unload(user.User, T) (any, error)
}

type RemoteAuther interface {
	Registration(model.User) (token string, err error)
	Authentication(model.User) (token string, err error)
}

type RemoteTexter[T any] interface {
	ReadAll[T]
	CRUD[T]
}

type RemoteByter[T any] interface {
	ReadAll[T]
	Loader[T]
	Read(user.User, T) (any, error)
	Delete(user.User, T) (any, error)
}
