package service

import (
	"github.com/boginskiy/GophKeeper/client/cmd/config"
	"github.com/boginskiy/GophKeeper/client/internal/api"
	"github.com/boginskiy/GophKeeper/client/internal/logg"
	"github.com/boginskiy/GophKeeper/client/internal/model"
	"github.com/boginskiy/GophKeeper/client/internal/user"
)

type TextService struct {
	Cfg          config.Config
	Logger       logg.Logger
	Type         string
	RemoteTexter api.RemoteTexter[model.Text]
}

func NewTextService(
	cfg config.Config,
	logger logg.Logger,
	remoteTexter api.RemoteTexter[model.Text]) *TextService {

	return &TextService{
		Cfg:          cfg,
		Logger:       logger,
		RemoteTexter: remoteTexter,
		Type:         "text",
	}
}

func (t *TextService) Create(user user.User, text model.Text) (any, error) {
	return t.RemoteTexter.Create(user, text)
}

func (t *TextService) Read(user user.User, text model.Text) (any, error) {
	return t.RemoteTexter.Read(user, text)
}

func (t *TextService) ReadAll(user user.User, text model.Text) (any, error) {
	return t.RemoteTexter.ReadAll(user, text)
}

func (t *TextService) Update(user user.User, text model.Text) (any, error) {
	return t.RemoteTexter.Update(user, text)
}

func (t *TextService) Delete(user user.User, text model.Text) (any, error) {
	return t.RemoteTexter.Delete(user, text)
}
