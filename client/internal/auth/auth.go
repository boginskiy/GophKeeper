package auth

import (
	"github.com/boginskiy/GophKeeper/client/cmd/client"
	"github.com/boginskiy/GophKeeper/client/cmd/config"
	"github.com/boginskiy/GophKeeper/client/internal/api"
	"github.com/boginskiy/GophKeeper/client/internal/cli"
	"github.com/boginskiy/GophKeeper/client/internal/errs"
	"github.com/boginskiy/GophKeeper/client/internal/logg"
	"github.com/boginskiy/GophKeeper/client/internal/model"
	"github.com/boginskiy/GophKeeper/client/internal/user"
	"github.com/boginskiy/GophKeeper/client/internal/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	ok, info := a.ErrorHandler(err)
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

	DecorGetEmail := a.Dialoger.GetEmailWithCheck(a.Dialoger.GetEmail, a.Dialoger.CheckEmail)
	_, err := DecorGetEmail(client, user)
	if err != nil {
		a.Logg.RaiseError(err, "bad email in authentication", nil)
		return false
	}

	DecorGetPassword := a.Dialoger.GetPasswordWithCheck(a.Dialoger.GetPassword, a.Dialoger.CheckPassword)
	_, err = DecorGetPassword(client, user)
	if err != nil {
		a.Logg.RaiseError(err, "bad password in authentication", nil)
		return false
	}

	a.Dialoger.ShowSomeInfo(client, "Authentication is successful")
	return true
}

func (a *Auth) ErrorHandler(err error) (bool, string) {
	if err == nil {
		return false, ""
	}

	// Ошибка локальная. Сервер не отвечает.
	if err == errs.ErrResponseServer {
		return true, "Server is unavailable, please try again later"
	}

	switch a.modifyErrServerOnCode(err) {
	// Ошибка создания пользователя.
	case codes.InvalidArgument:
		return true, "User creation error"

	// Ошибка уникального email.
	case codes.AlreadyExists:
		return true, "Unique email error"

	// Ошибка создания токена или когда кладем токен в заголовок.
	case codes.Internal:
		return true, "An error in creating or transferring a token"

	// Неизвестная ошибка.
	default:
		return true, "Unknown error"
	}
}

func (a *Auth) modifyErrServerOnCode(err error) codes.Code {
	statusErr, ok := status.FromError(err)
	if !ok {
		return codes.OK
	}
	return statusErr.Code()
}
