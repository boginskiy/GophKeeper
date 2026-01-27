package api

import (
	"github.com/boginskiy/GophKeeper/client/internal/model"
	"github.com/boginskiy/GophKeeper/client/internal/user"
)

type Auther interface {
	Registration(model.User) (token string, err error)
	Authentication(user model.User) (token string, err error)
}

type CRUD interface {
	Create(user *user.UserCLI, text model.Text) (any, error)
	Read(user *user.UserCLI, text model.Text) (any, error)
	ReadAll(user *user.UserCLI, text model.Text) (any, error)
	Update(user *user.UserCLI, text model.Text) (any, error)
	// Delete(user *user.UserCLI, text model.Text) (any, error)
}

type ServiceAPI interface {
	Auther
	CRUD
}
