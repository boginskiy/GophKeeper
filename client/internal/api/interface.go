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
	// Read(user *user.UserCLI, text model.Text)
	// ReadAll(user *user.UserCLI, text model.Text)
	// Update(user *user.UserCLI, text model.Text)
	// Delete(user *user.UserCLI, text model.Text)
}

type ServiceAPI interface {
	Auther
	CRUD
}
