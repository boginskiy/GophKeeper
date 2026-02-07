package infra

import (
	"os"

	"github.com/boginskiy/GophKeeper/client/internal/model"
	"github.com/boginskiy/GophKeeper/client/internal/user"
)

// DataRecover
type DataRecover interface {
	RecoveryPassword(email string) bool
	GetRandomNumber() string
}

// Checker
type Checker interface {
	CheckTwoString(oneStr, twoStr string) bool
	CheckPassword(userPassword, password string) bool
	CheckTypeFile(pathToFile, typ string) bool
}

// Dialoger
type Dialoger interface {
	DataGetter
	Verifer
	Getter
	Shower

	CallWindows(string) (string, error)
}

// FileManager
type FileManager interface {
	CreateFileInStore(*model.Bytes) (file *os.File, path string, err error)
	GetModelBytesFromFile(pathToFile string, typ string) (*model.Bytes, error)
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
