package user

import (
	"bufio"
	"context"
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
	Name     string
	InMess   chan string
	OutMess  chan string
	Scanner  *bufio.Scanner
	Logger   logg.Logger
	User     *model.User
	Identity *Identity
}

func NewUserCLI(ctx context.Context, logger logg.Logger, out, in chan string, identity *Identity) *UserCLI {
	tmp := &UserCLI{
		Name:     NAMECLI,
		OutMess:  out,
		InMess:   in,
		Logger:   logger,
		Scanner:  bufio.NewScanner(os.Stdin),
		Identity: identity,
	}

	// Identification user
	identity.SystemIdentification(tmp)

	return tmp
}

func (u *UserCLI) SaveConfig() {
	err := u.Identity.SaveCurrentUser(u.User)
	if err != nil {
		u.Logger.RaiseError(err, "error saving new user i config file", nil)
	}
}

func (u *UserCLI) TakeSystemInfoCurrentUser() (username, uid string) {
	user, err := user.Current()
	if err != nil {
		u.Logger.RaiseError(err, "error taking extra user info", nil)
		return
	}
	return user.Username, user.Uid
}

func (u *UserCLI) ReceiveMess() (string, error) {
	time.Sleep(500 * time.Millisecond)
	fmt.Fprintf(os.Stdout, "%s: ", u.Name)

	if !u.Scanner.Scan() {
		return "", errors.New("error user reciver")
	}
	return u.Scanner.Text(), nil
}

func (u *UserCLI) NewUser(userName, email, phone, password string) {
	// Save system info about new user
	systemName, systemId := u.TakeSystemInfoCurrentUser()
	// Hash password
	hash, err := pkg.GenerateHash(password)
	u.Logger.CheckWithFatal(err, "error in hashing password")

	tmp := &model.User{
		UserName:       userName,
		Email:          email,
		PhoneNumber:    phone,
		Password:       hash,
		SystemUserName: systemName,
		SystemUserId:   systemId,
	}
	u.User = tmp
}
