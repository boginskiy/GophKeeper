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
	CreateText(user *user.UserCLI, text model.Text)
	ReadText(user *user.UserCLI, text model.Text)
	UpdateText(user *user.UserCLI, text model.Text)
}

type ServiceAPI interface {
	Auther
	CRUD
}
