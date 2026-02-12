package auth

import (
	"github.com/boginskiy/GophKeeper/client/cmd/config"
	"github.com/boginskiy/GophKeeper/client/internal/api"
	"github.com/boginskiy/GophKeeper/client/internal/logg"
	"github.com/boginskiy/GophKeeper/client/internal/model"
	"github.com/boginskiy/GophKeeper/client/internal/user"
)

type AuthService struct {
	Cfg          config.Config
	Logger       logg.Logger
	Identity     Identifier
	RemoteAuther api.RemoteAuther
}

func NewAuthService(
	config config.Config,
	logger logg.Logger,
	identity Identifier,
	remoteAuther api.RemoteAuther,
) *AuthService {

	return &AuthService{
		Cfg:          config,
		Logger:       logger,
		Identity:     identity,
		RemoteAuther: remoteAuther,
	}
}

func (a *AuthService) Identification(user user.User) bool {
	return a.Identity.Identification(user)
}

func (a *AuthService) Registration(user user.User, modUser *model.User) (string, error) {
	token, err := a.RemoteAuther.Registration(*modUser)

	// Обработка ошибок
	ok, info := ErrorHandler(err)
	if ok {
		return info, err
	}

	user.SaveLocalUser(modUser)
	user.GetModelUser().Token = token
	return token, nil
}

func (a *AuthService) Authentication(user user.User, modUser *model.User) (string, error) {
	token, err := a.RemoteAuther.Authentication(*modUser)

	// Обработка ошибок
	ok, info := ErrorHandler(err)
	if ok {
		return info, err
	}

	user.GetModelUser().Token = token
	return token, nil
}

func (a *AuthService) Recovery(user user.User, modUser *model.User) (string, error) {
	token, err := a.RemoteAuther.Recovery(*modUser)

	// Обработка ошибок
	ok, info := ErrorHandler(err)
	if ok {
		return info, err
	}

	user.SaveLocalUser(modUser)
	user.GetModelUser().Token = token
	return token, nil
}
