package auth

import (
	"fmt"

	"github.com/boginskiy/GophKeeper/client/cmd/client"
	"github.com/boginskiy/GophKeeper/client/cmd/config"
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
}

func NewAuth(
	config config.Config,
	logger logg.Logger,
	userch chan *model.User,
	fileHdlr utils.FileHandler,
	identity *Identity,
	dialoger cli.Dialoger,
) *Auth {

	return &Auth{
		Cfg:         config,
		Logg:        logger,
		UserChan:    userch,
		FileHendler: fileHdlr,
		Identity:    identity,
		Dialoger:    dialoger,
	}
}

func (a *Auth) Registration(isThereAuthent bool, client *client.ClientCLI, user *user.UserCLI) bool {
	if isThereAuthent {
		return false
	}

	a.Dialoger.ShowRegister(client, user)

	userName, err := a.Dialoger.GetUserName(client, user)
	a.Logg.CheckWithFatal(err, "bad user name")

	email, err := a.Dialoger.GetEmail(client, user)
	a.Logg.CheckWithFatal(err, "bad email")

	phone, err := a.Dialoger.GetPhone(client, user)
	a.Logg.CheckWithFatal(err, "bad phone")

	password, err := a.Dialoger.GetPassword(client, user)
	a.Logg.CheckWithFatal(err, "bad password")

	// Вот как это разрулить ?
	// Доделать и написать тесты
	//

	// Передача данных на сервис AuthRemote для регистрации на удаленном сервере.
	a.UserChan <- model.NewUser(userName, email, phone, password)
	newUser := <-a.UserChan

	fmt.Printf("%+v\n", newUser)

	// ERROR ...
	if newUser.StatusError != nil {

		// Ошибка локальная. Сервер не отвечает.
		if newUser.StatusError == errs.ErrResponseServer {
			client.SendMess("Server is unavailable, please try again later")
			return false
		}

		// Ошибки с сервера.
		codeErr := a.CodeErrFromServerGRPC(newUser.StatusError)

		// Ошибка создания пользователя.
		if codeErr == codes.InvalidArgument {
			// TODO...
		}

		// Ошибка уникального email.
		if codeErr == codes.AlreadyExists {
			// TODO...
		}

		// Ошибка создания токена.
		// Ошибка когда кладем токен в заголовок.
		if codeErr == codes.Internal {
			// TODO...
		}

	}

	user.PreparUser(newUser)
	client.SendMess("Registration is successful")
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

	return true
}

func (a *Auth) CodeErrFromServerGRPC(err error) codes.Code {
	statusErr, ok := status.FromError(err)
	if !ok {
		return codes.OK
	}
	return statusErr.Code()
}
