package auth

import (
	"fmt"

	"github.com/boginskiy/GophKeeper/client/internal/client"
	"github.com/boginskiy/GophKeeper/client/internal/logg"
	"github.com/boginskiy/GophKeeper/client/internal/user"
	"github.com/boginskiy/GophKeeper/client/pkg"
)

const ATTEMPTS = 3

type Auth struct {
	Logger logg.Logger
}

func NewAuth(logger logg.Logger) *Auth {
	return &Auth{Logger: logger}
}

func (a *Auth) Authorization(client *client.ClientCLI, user *user.UserCLI) {
	a.Welcome(client, user)

	// Check of registration.
	IsRegistration := user.User != nil

	// Check of authentication.
	IsAuthentication := IsRegistration && a.Authentication(client, user)

	fmt.Println(IsAuthentication)

	// New registration.
	_ = !IsAuthentication && a.Registration(client, user)

}

func (a *Auth) Registration(client *client.ClientCLI, user *user.UserCLI) bool {
	// user name
	client.SendMess(
		"You need to register...",
		"Enter user name...")
	userName, err := user.ReceiveMess() // TODO нужна проверка
	a.Logger.CheckWithFatal(err, "bad user name")

	// email
	email, err := GetEmail(client, user) // TODO нужна проверка c запросом на сервер
	a.Logger.CheckWithFatal(err, "bad email")

	// phone
	client.SendMess("Enter phone...")
	phone, err := user.ReceiveMess() // TODO нужна проверка
	a.Logger.CheckWithFatal(err, "bad phone")

	// password
	password, err := GetPassword(client, user) // TODO нужна проверка
	a.Logger.CheckWithFatal(err, "bad password")

	user.NewUser(userName, email, phone, password)
	client.SendMess("Registration is successful")
	return true
}

func (a *Auth) Authentication(client *client.ClientCLI, user *user.UserCLI) bool {
	client.SendMess("You need log in...")

	DecorGetEmail := TryToGetSeveralTimes(GetEmail, a.checkEmail)
	_, err := DecorGetEmail(client, user)
	if err != nil {
		a.Logger.RaiseError(err, "bad email in authentication", nil)
		return false
	}

	DecorGetPassword := TryToGetSeveralTimes(GetPassword, a.checkPassword)
	_, err = DecorGetPassword(client, user)
	if err != nil {
		a.Logger.RaiseError(err, "bad password in authentication", nil)
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

func (a *Auth) Welcome(client *client.ClientCLI, user *user.UserCLI) {
	client.SendMess("Hello! Press the 'Enter'...")
	user.ReceiveMess()
}
