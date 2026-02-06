package comm

import (
	"fmt"

	"github.com/boginskiy/GophKeeper/client/cmd/client"
	"github.com/boginskiy/GophKeeper/client/internal/auth"
	"github.com/boginskiy/GophKeeper/client/internal/errs"
	"github.com/boginskiy/GophKeeper/client/internal/infra"
	"github.com/boginskiy/GophKeeper/client/internal/model"
	"github.com/boginskiy/GophKeeper/client/internal/user"
)

type Root struct {
	Dialoger  infra.Dialoger
	CommText  Commander
	CommMedia Commander
}

func NewRoot(
	dialoger infra.Dialoger,
	commtext Commander,
	commmedia Commander,
) *Root {

	return &Root{
		Dialoger:  dialoger,
		CommText:  commtext,
		CommMedia: commmedia,
	}
}

func (r *Root) ExecuteComm(authOK bool, client *client.ClientCLI, user *user.UserCLI) {
authLoop:
	for authOK {
		// Get command.
		comm, _ := r.Dialoger.GetSomeThing(

			fmt.Sprintf("%s\n\r%s",
				"Enter the data type you want to work with: \n\r\t bytes \n\r\t text \n\r\t image  \n\r\t sound \n\r\t video",
				"go out: exit, need help: help"))

		switch comm {
		case "exit", "help":
			break authLoop

		case "bytes":
			r.CommMedia.Registration(user, "bytes")
		case "text":
			r.CommText.Registration(user, "text")
		case "image":
			r.CommMedia.Registration(user, "image")
		case "sound":
			r.CommMedia.Registration(user, "sound")
		case "video":
			r.CommMedia.Registration(user, "video")

		default:
			r.Dialoger.ShowIt("Invalid command. Try again...")
		}
	}
}

func (r *Root) ExecuteAuth(authSrv auth.Auth, user user.User) bool {
	var info string

	// Identification.
	if ok := authSrv.Identification(user); ok {

		// Authentication.
		email, password, err := r.Dialoger.VerifyDataAuth(user)
		if err == nil {
			checkUser := &model.User{Email: email, Password: password}
			info, err = authSrv.Authentication(user, checkUser)
			if err == nil {
				r.Dialoger.ShowIt("Authentication is successful")
				return true
			}
		}

		// Пользователь забыл password. Идем восстанавливать.
		if err == errs.ErrPassword {

		}

		// Пользователь забыл email. Идем восстанавливать.
		if err == errs.ErrEmail {

		}
	}

	// Registration.
	r.Dialoger.ShowIt(info)

	newUser := model.NewUser(r.Dialoger.GetDataRegister())
	info, err := authSrv.Registration(user, newUser)
	if err == nil {
		r.Dialoger.ShowIt("Registration is successful")
		return true
	}
	r.Dialoger.ShowIt(info)
	return false
}
