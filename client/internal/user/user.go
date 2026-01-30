package user

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/user"
	"time"

	"github.com/boginskiy/GophKeeper/client/internal/logg"
	"github.com/boginskiy/GophKeeper/client/internal/model"
	"github.com/boginskiy/GophKeeper/client/pkg"
)

const NAMECLI = "USER"

type UserCLI struct {
	Name    string
	Scanner *bufio.Scanner
	Logg    logg.Logger
	User    *model.User
}

func NewUserCLI(logger logg.Logger) *UserCLI {
	return &UserCLI{
		Name:    NAMECLI,
		Logg:    logger,
		Scanner: bufio.NewScanner(os.Stdin),
	}
}

func (u *UserCLI) GetModelUser() *model.User {
	return u.User
}

func (u *UserCLI) GetSystemInfo() (username, uid string) {
	user, err := user.Current()
	if err != nil {
		u.Logg.RaiseError(err, "error taking extra user info", nil)
		return
	}
	return user.Username, user.Uid
}

func (u *UserCLI) ReceiveMess() (string, error) {
	time.Sleep(200 * time.Millisecond)
	fmt.Fprintf(os.Stdout, "%s: ", u.Name)

	if !u.Scanner.Scan() {
		return "", errors.New("error user reciver")
	}
	return u.Scanner.Text(), nil
}

func (u *UserCLI) SaveLocalUser(localUser *model.User) {
	systemName, systemId := u.GetSystemInfo()         // Save system info about new user
	hash, err := pkg.GenerateHash(localUser.Password) // Hash password

	u.Logg.CheckWithFatal(err, "error in hashing password")

	localUser.SystemUserName = systemName
	localUser.SystemUserId = systemId
	localUser.Password = hash

	u.User = localUser
}

func (u *UserCLI) SavePreviosUser(previosUser *model.User) {
	u.User = previosUser
}
