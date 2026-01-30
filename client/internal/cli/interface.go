package cli

import (
	"github.com/boginskiy/GophKeeper/client/internal/user"
)

type Getter interface {
	GetEnterIt(string) (string, error)
	GetSomeThing(string) (string, error)
}

type DataGetter interface {
	GetDataAction(action string) string
	GetDataRegister() (userName, email, phone, password string)
}

type Shower interface {
	ShowIt(string)
	ShowErr(error)
}

type Verifer interface {
	VerifyEnterPassword(needToTake, needToCompare string, quantity int) bool
	VerifyEnterIt(needToTake, needToCompare string, quantity int) bool
	VerifyDataAuth(user.User) bool
}

type ShowGetter interface {
	DataGetter
	Verifer
	Getter
	Shower
}
