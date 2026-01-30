package api

import (
	"github.com/boginskiy/GophKeeper/client/internal/model"
	"github.com/boginskiy/GophKeeper/client/internal/user"
)

type Auther interface {
	Registration(model.User) (token string, err error)
	Authentication(model.User) (token string, err error)
}

type CRUD interface {
	Create(user.User, model.Text) (any, error)
	Read(user.User, model.Text) (any, error)
	ReadAll(user.User, model.Text) (any, error)
	Update(user.User, model.Text) (any, error)
	Delete(user.User, model.Text) (any, error)
}

type Uploader interface {
	Upload(*user.UserCLI, model.Bytes) (any, error)
}

type ServiceAPI interface {
	Uploader
	Auther
	CRUD
}
