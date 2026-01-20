package auth

import (
	"fmt"

	"github.com/boginskiy/GophKeeper/client/cmd/client"
	"github.com/boginskiy/GophKeeper/client/cmd/config"
	"github.com/boginskiy/GophKeeper/client/internal/logg"
	"github.com/boginskiy/GophKeeper/client/internal/model"
	"github.com/boginskiy/GophKeeper/client/internal/user"
	"github.com/boginskiy/GophKeeper/client/internal/utils"
	"github.com/boginskiy/GophKeeper/client/pkg"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth struct {
	Cfg         config.Config
	Logg        logg.Logger
	UserChan    chan *model.User
	FileHendler utils.FileHandler
	Identity    *Identity
}

func NewAuth(
	config config.Config,
	logger logg.Logger,
	userch chan *model.User,
	fileHdlr utils.FileHandler,
	identity *Identity) *Auth {

	return &Auth{
		Cfg:         config,
		Logg:        logger,
		UserChan:    userch,
		FileHendler: fileHdlr,
		Identity:    identity}
}

func (a *Auth) Identification(client *client.ClientCLI, user *user.UserCLI) {
	// Identification user. load data from config.file if exist.
	ok := a.Identity.Identification(user)
	if ok {
		// TODO!
		// После успешной идентификации пользователя, подгружаются данные из
		// config.file и после можно обращаться к нему по имени.
	}

	// Check of registration.
	IsRegistration := user.User != nil

	// Check of authentication.
	IsAuthentication := IsRegistration && a.Authentication(client, user)

	// New registration.
	_ = !IsAuthentication && a.Registration(client, user)

}

func (a *Auth) Registration(client *client.ClientCLI, user *user.UserCLI) bool {
	// user name
	client.SendMess(
		"You need to register...",
		"Enter user name...")
	userName, err := user.ReceiveMess() // TODO нужна проверка
	a.Logg.CheckWithFatal(err, "bad user name")

	// email
	email, err := GetEmail(client, user) // TODO нужна проверка c запросом на сервер
	a.Logg.CheckWithFatal(err, "bad email")

	// phone
	client.SendMess("Enter phone...")
	phone, err := user.ReceiveMess() // TODO нужна проверка
	a.Logg.CheckWithFatal(err, "bad phone")

	// password
	password, err := GetPassword(client, user) // TODO нужна проверка
	a.Logg.CheckWithFatal(err, "bad password")

	// Передача данных на сервис AuthRemote для регистрации на удаленном сервере.
	a.UserChan <- model.NewUser(userName, email, phone, password)
	newUser := <-a.UserChan

	fmt.Printf("%+v\n", newUser)

	// // ERROR ...
	// if newUser.StatusError != nil {

	// 	// Ошибка локальная. Сервер не отвечает.
	// 	if newUser.StatusError == errs.ErrResponseServer {
	// 		client.SendMess("Server is unavailable, please try again later")
	// 		return false
	// 	}

	// 	// Ошибки с сервера.
	// 	codeErr := a.CodeErrFromServerGRPC(newUser.StatusError)

	// 	// Ошибка создания пользователя.
	// 	if codeErr == codes.InvalidArgument {
	// 		// TODO...
	// 	}

	// 	// Ошибка уникального email.
	// 	if codeErr == codes.AlreadyExists {
	// 		// TODO...
	// 	}

	// 	// Ошибка создания токена.
	// 	// Ошибка когда кладем токен в заголовок.
	// 	if codeErr == codes.Internal {
	// 		// TODO...
	// 	}

	// }

	// user.PreparUser(newUser)
	client.SendMess("Registration is successful")
	return true
}

func (a *Auth) CodeErrFromServerGRPC(err error) codes.Code {
	statusErr, ok := status.FromError(err)
	if !ok {
		return codes.OK
	}
	return statusErr.Code()
}

// if err == errs.ErrCreateUser {
// 	// Ошибка создания пользователя.
// 	return nil, status.Errorf(codes.InvalidArgument, "%s", err.Error())
// }

// if err == errs.ErrEmailNotUnique {
// 	// Ошибка уникального email.
// 	return nil, status.Errorf(codes.AlreadyExists, "%s", err.Error())
// }

// if err == errs.ErrCreateToken {
// 	// Ошибка создания токена.
// 	return nil, status.Errorf(codes.Internal, "%s", err.Error())
// }

// // Кладем токен в заголовок
// err = k.putDataToCtx(ctx, "authorization", token)
// if err != nil {
// 	return nil, status.Errorf(codes.Internal, "%s", err.Error())

func (a *Auth) Authentication(client *client.ClientCLI, user *user.UserCLI) bool {
	client.SendMess("You need log in...")

	DecorGetEmail := TryToGetSeveralTimes(GetEmail, a.checkEmail)
	_, err := DecorGetEmail(client, user)
	if err != nil {
		a.Logg.RaiseError(err, "bad email in authentication", nil)
		return false
	}

	DecorGetPassword := TryToGetSeveralTimes(GetPassword, a.checkPassword)
	_, err = DecorGetPassword(client, user)
	if err != nil {
		a.Logg.RaiseError(err, "bad password in authentication", nil)
		return false
	}

	client.SendMess("Authentication is successful")
	return true
}

func (a *Auth) checkEmail(user *user.UserCLI, email string) bool {
	return user.User.Email == email
}

func (a *Auth) checkPassword(user *user.UserCLI, password string) bool {
	return pkg.CompareHashAndPassword(user.User.Password, password)
}
