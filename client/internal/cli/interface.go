package cli

import (
	"github.com/boginskiy/GophKeeper/client/cmd/client"
	"github.com/boginskiy/GophKeeper/client/internal/user"
)

// TODO! Это вынести в отдельный сервис. Пока так оставил
type Checker interface {
	CheckPassword(userPassword, password string) bool
	CheckEmail(userEmail, email string) bool
}

type GetterEmail interface {
	GetEmailWithCheck(GetterFn, CheckerFn) GetterFn
	GetEmail(*client.ClientCLI, *user.UserCLI) (string, error)
}

type GetterPassword interface {
	GetPasswordWithCheck(GetterFn, CheckerFn) GetterFn
	GetPassword(*client.ClientCLI, *user.UserCLI) (string, error)
}

type GetterUserName interface {
	GetUserName(*client.ClientCLI, *user.UserCLI) (string, error)
}

type Getter interface {
	// GetCommand(*client.ClientCLI, *user.UserCLI) (string, error)
	GetSomeThing(*client.ClientCLI, *user.UserCLI, string) (string, error)
}

type GetterPhone interface {
	GetPhone(*client.ClientCLI, *user.UserCLI) (string, error)
}

type Dialoger interface {
	GetterUserName
	GetterPassword
	Getter
	GetterPhone
	GetterEmail
	Checker

	ShowStatusAuth(*client.ClientCLI, *user.UserCLI)
	ShowRegister(*client.ClientCLI, *user.UserCLI)
	ShowHello(*client.ClientCLI, *user.UserCLI)
	ShowLogIn(*client.ClientCLI, *user.UserCLI)
	ShowSomeInfo(*client.ClientCLI, string)

	DialogsAbRegister(*client.ClientCLI, *user.UserCLI) (userName, email, phone, password string)
}
