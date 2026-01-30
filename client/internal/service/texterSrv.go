package service

import (
	"github.com/boginskiy/GophKeeper/client/cmd/config"
	"github.com/boginskiy/GophKeeper/client/internal/api"
	"github.com/boginskiy/GophKeeper/client/internal/logg"
	"github.com/boginskiy/GophKeeper/client/internal/model"
	"github.com/boginskiy/GophKeeper/client/internal/user"
)

type TexterService struct {
	Cfg        config.Config
	Logg       logg.Logger
	ServiceAPI api.ServiceAPI
	Type       string
}

func NewTexterService(
	cfg config.Config,
	logger logg.Logger,
	serviceAPI api.ServiceAPI,
) *TexterService {

	return &TexterService{
		Cfg:        cfg,
		Logg:       logger,
		ServiceAPI: serviceAPI,
		Type:       "text",
	}
}

func (t *TexterService) Create(user user.User, text model.Text) (any, error) {
	return t.ServiceAPI.Create(user, text)
}

func (t *TexterService) Read(user user.User, text model.Text) (any, error) {
	return t.ServiceAPI.Read(user, text)
}

func (t *TexterService) ReadAll(user user.User, text model.Text) (any, error) {
	return t.ServiceAPI.ReadAll(user, text)
}

func (t *TexterService) Update(user user.User, text model.Text) (any, error) {
	return t.ServiceAPI.Update(user, text)
}

func (t *TexterService) Delete(user user.User, text model.Text) (any, error) {
	return t.ServiceAPI.Delete(user, text)
}
