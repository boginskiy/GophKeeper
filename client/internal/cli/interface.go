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

type Getter interface {
	GetIt(*client.ClientCLI, *user.UserCLI, string) (string, error)
	GetSomeThing(*client.ClientCLI, *user.UserCLI, string) (string, error)
}

type Shower interface {
	ShowHello(*client.ClientCLI, *user.UserCLI)
	ShowIt(*client.ClientCLI, string)
}

type Dialoger interface {
	GetterPassword
	GetterEmail
	Checker
	Getter
	Shower

	DialogsAbRegister(*client.ClientCLI, *user.UserCLI) (userName, email, phone, password string)
	DialogsAbAction(*client.ClientCLI, *user.UserCLI, string) string
}
