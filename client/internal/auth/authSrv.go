package auth

import (
	"github.com/boginskiy/GophKeeper/client/cmd/config"
	"github.com/boginskiy/GophKeeper/client/internal/api"
	"github.com/boginskiy/GophKeeper/client/internal/errs"
	"github.com/boginskiy/GophKeeper/client/internal/logg"
	"github.com/boginskiy/GophKeeper/client/internal/model"
	"github.com/boginskiy/GophKeeper/client/internal/user"
	"github.com/boginskiy/GophKeeper/client/internal/utils"
)

type AuthService struct {
	Cfg         config.Config
	Logg        logg.Logger
	FileHendler utils.FileHandler
	Identity    Identifier
	// DialogSrv   cli.ShowGetter
	ServiceAPI api.ServiceAPI
}

func NewAuthService(
	config config.Config,
	logger logg.Logger,
	fileHdlr utils.FileHandler,
	identity Identifier,
	// dialog cli.ShowGetter,
	serviceAPI api.ServiceAPI,
) *AuthService {

	return &AuthService{
		Cfg:         config,
		Logg:        logger,
		FileHendler: fileHdlr,
		Identity:    identity,
		// DialogSrv:   dialog,
		ServiceAPI: serviceAPI,
	}
}

func (a *AuthService) Identification(user user.User) bool {
	return a.Identity.Identification(user)
}

func (a *AuthService) Registration(user user.User, modUser *model.User) (string, error) {
	token, err := a.ServiceAPI.Registration(*modUser)

	// Обработка ошибок
	ok, info := ErrorHandler(err)
	if ok {
		return info, err
	}

	user.SaveLocalUser(modUser)
	user.GetModelUser().Token = token
	return token, nil
}

func (a *AuthService) Authentication(verify bool, user user.User) (string, error) {
	if !verify {
		return "uncorrected credentials", errs.ErrUncorrectCredentials
	}
	token, err := a.ServiceAPI.Authentication(*user.GetModelUser())

	// Обработка ошибок
	ok, info := ErrorHandler(err)
	if ok {
		return info, err
	}

	user.GetModelUser().Token = token
	return token, nil
}
