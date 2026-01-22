package auth

import (
	"github.com/boginskiy/GophKeeper/client/cmd/client"
	"github.com/boginskiy/GophKeeper/client/cmd/config"
	"github.com/boginskiy/GophKeeper/client/internal/api"
	"github.com/boginskiy/GophKeeper/client/internal/cli"
	"github.com/boginskiy/GophKeeper/client/internal/logg"
	"github.com/boginskiy/GophKeeper/client/internal/model"
	"github.com/boginskiy/GophKeeper/client/internal/user"
	"github.com/boginskiy/GophKeeper/client/internal/utils"
)

type Auth struct {
	Cfg         config.Config
	Logg        logg.Logger
	UserChan    chan *model.User
	FileHendler utils.FileHandler
	Identity    *Identity
	Dialoger    cli.Dialoger
	ServiceAPI  api.ServiceAPI
}

func NewAuth(
	config config.Config,
	logger logg.Logger,
	userch chan *model.User,
	fileHdlr utils.FileHandler,
	identity *Identity,
	dialoger cli.Dialoger,
	serviceAPI api.ServiceAPI,
) *Auth {

	return &Auth{
		Cfg:         config,
		Logg:        logger,
		UserChan:    userch,
		FileHendler: fileHdlr,
		Identity:    identity,
		Dialoger:    dialoger,
		ServiceAPI:  serviceAPI,
	}
}

func (a *Auth) Registration(isThereAuthent bool, client *client.ClientCLI, user *user.UserCLI) bool {
	if isThereAuthent {
		return false
	}

	// Передача данных на сервис ServiceAPI для регистрации на удаленном сервере.
	userName, email, phone, password := a.Dialoger.DialogsAbRegister(client, user)
	newUser := model.NewUser(userName, email, phone, password)
	token, err := a.ServiceAPI.Registration(*newUser)

	// Обработка ошибок
	ok, info := ErrorHandler(err)
	if ok {
		a.Dialoger.ShowSomeInfo(client, info)
		return false
	}

	newUser.Token = token
	newUser.StatusError = err

	user.SaveLocalUser(newUser)
	a.Dialoger.ShowSomeInfo(client, "Registration is successful")
	return true
}

func (a *Auth) Authentication(isThereRegistr bool, client *client.ClientCLI, user *user.UserCLI) bool {
	if !isThereRegistr {
		return false
	}

	a.Dialoger.ShowLogIn(client, user)

	// Работает так: вводишь с CLI данные email
	// функция осуществляет сверку введенного email  с email из config
	// если нет совпадения, то есть 3 попытки ввести верный email,
	// иначе аутентификация признается невалидной.

	// TODO! Это можно реализовать в интерфейсе Dialoger. Убрать отсюда.
	DecorGetEmail := a.Dialoger.GetEmailWithCheck(a.Dialoger.GetEmail, a.Dialoger.CheckEmail)
	email, err := DecorGetEmail(client, user)
	if err != nil {
		a.Logg.RaiseError(err, "bad email in authentication", nil)
		return false
	}

	// TODO! Это можно реализовать в интерфейсе Dialoger. Убрать отсюда.
	DecorGetPassword := a.Dialoger.GetPasswordWithCheck(a.Dialoger.GetPassword, a.Dialoger.CheckPassword)
	password, err := DecorGetPassword(client, user)
	if err != nil {
		a.Logg.RaiseError(err, "bad password in authentication", nil)
		return false
	}

	token, err := a.ServiceAPI.Authentication(model.User{Email: email, Password: password})

	// Обработка ошибок
	ok, info := ErrorHandler(err)
	if ok {
		a.Dialoger.ShowSomeInfo(client, info)
		return false
	}

	user.User.Token = token
	user.User.StatusError = err

	a.Dialoger.ShowSomeInfo(client, "Authentication is successful")
	return true
}
