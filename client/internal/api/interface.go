package api

import "github.com/boginskiy/GophKeeper/client/internal/model"

type ServiceAPI interface {
	Registration(model.User) (token string, err error)
	Authentication(user model.User) (token string, err error)
}
