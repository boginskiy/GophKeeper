package comm

import (
	"context"
	"fmt"

	"github.com/boginskiy/GophKeeper/client/cmd/client"
	"github.com/boginskiy/GophKeeper/client/internal/auth"
	"github.com/boginskiy/GophKeeper/client/internal/errs"
	"github.com/boginskiy/GophKeeper/client/internal/infra"
	"github.com/boginskiy/GophKeeper/client/internal/model"
	"github.com/boginskiy/GophKeeper/client/internal/user"
)

type Root struct {
	Ctx      context.Context
	Dialoger infra.Dialoger
	// DataRecover infra.DataRecover
	CommText  Commander
	CommMedia Commander
	MailChan  chan<- string
	CodeChan  <-chan string
}

func NewRoot(
	ctx context.Context,
	dialoger infra.Dialoger,
	// dataRecover infra.DataRecover,
	commtext Commander,
	commmedia Commander,
	mailChan chan<- string,
	codeChan <-chan string,
) *Root {

	return &Root{
		Ctx:      ctx,
		Dialoger: dialoger,
		// DataRecover: dataRecover,
		CommText:  commtext,
		CommMedia: commmedia,
		MailChan:  mailChan,
		CodeChan:  codeChan,
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
			r.Dialoger.ShowIt(info)
		}

		// Recovery Password.
		if err == errs.ErrPassword {
			comm, _ := r.Dialoger.GetSomeThing(
				fmt.Sprintf("%s\n\r%s",
					"Do you want to recover your password? \n\r\t yes \n\r\t no",
					"go out: exit, need help: help"))

			if comm == "yes" {
				// Отправляем на обработку в сервис восстановления пароля.
				r.MailChan <- email

				codeUser, _ := r.Dialoger.GetSomeThing(
					fmt.Sprintf("%s",
						"Enter code to recover from your mail..."))

				select {
				case <-r.Ctx.Done():
					r.Dialoger.ShowIt("Some problems with data recovery")

				case codeGener := <-r.CodeChan:
					// Check user's code with generation code
					if codeUser != codeGener {
						r.Dialoger.ShowIt("Invalid code to recover")

					} else {
						newPassword, _ := r.Dialoger.GetSomeThing(
							fmt.Sprintf("%s",
								"Enter a new password..."))

						updateUser := &model.User{Password: newPassword, Email: email}
						info, err := authSrv.Recovery(user, updateUser)

						if err == nil {
							r.Dialoger.ShowIt("Recovery password is successful")
							return true
						}
						r.Dialoger.ShowIt(info)
						return false
					}
				}
			}
		}
	}

	// Registration.
	newUser := model.NewUser(r.Dialoger.GetDataRegister())
	info, err := authSrv.Registration(user, newUser)
	if err == nil {
		r.Dialoger.ShowIt("Registration is successful")
		return true
	}
	r.Dialoger.ShowIt(info)
	return false
}

// TODO можно сделать отдельной страницей как
// Media ?
