package infra

import "github.com/boginskiy/GophKeeper/client/internal/user"

// Checker for Check
type Checker interface {
	CheckTwoString(oneStr, twoStr string) bool
	CheckPassword(userPassword, password string) bool
	CheckTypeFile(pathToFile, typ string) bool
}

// ShowGetter for DialogService
type ShowGetter interface {
	DataGetter
	Verifer
	Getter
	Shower

	CallWindows(string) (string, error)
}

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
	VerifyEnterPassword(needToTake, needToCompare string, quantity int) (string, error)
	VerifyEnterIt(needToTake, needToCompare string, quantity int) (string, error)
	VerifyDataAuth(user.User) (email, password string, err error)
}
