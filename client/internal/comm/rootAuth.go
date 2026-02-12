package comm

import (
	"context"
	"fmt"
	"time"

	"github.com/boginskiy/GophKeeper/client/internal/auth"
	"github.com/boginskiy/GophKeeper/client/internal/errs"
	"github.com/boginskiy/GophKeeper/client/internal/infra"
	"github.com/boginskiy/GophKeeper/client/internal/model"
	"github.com/boginskiy/GophKeeper/client/internal/user"
)

type RootAuth struct {
	Ctx         context.Context
	MailChan    chan<- string
	CodeChan    <-chan string
	Dialoger    infra.Dialoger
	AuthService auth.Auth
}

func NewRootAuth(
	ctx context.Context,
	mailChan chan<- string,
	codeChan <-chan string,
	dialoger infra.Dialoger,
	authService auth.Auth,
) *RootAuth {

	return &RootAuth{
		Ctx:         ctx,
		MailChan:    mailChan,
		CodeChan:    codeChan,
		Dialoger:    dialoger,
		AuthService: authService,
	}
}

func (r *RootAuth) ExecuteComm(_ bool, user user.User) bool {
	// Identification.
	identifIsOk := r.AuthService.Identification(user)

	// Authentication.
	authIsOk, err := r.ExecuteAuth(identifIsOk, user)
	if authIsOk {
		return true
	}

	// Recovery.
	recoverIsOk := r.ExecuteRecover(err, user)
	if recoverIsOk {
		return true
	}

	// Registration.
	return r.ExecuteRegistr(user)
}

func (r *RootAuth) ExecuteRecover(err error, user user.User) bool {
	if err != errs.ErrPassword {
		return false
	}

	comm, _ := r.Dialoger.GetSomeThing(
		fmt.Sprintf("%s\n\r%s",
			"Do you want to recover your password? \n\r\t yes \n\r\t no",
			"go out: exit, need help: help"))

	if comm == "no" {
		return false
	}

	// Send a mail to in RecoveryService
	r.MailChan <- user.GetModelUser().Email

	// Просим пользователя ввести полученный по почте code.
	codeUser, _ := r.Dialoger.GetSomeThing(
		fmt.Sprintf("%s",
			"Enter code to recover from your mail..."))

	// Geting one-time password from RecoveryService
	codeGener := r.getAnswerFromRecoveryService(r.Ctx)

	if codeUser != codeGener {
		r.Dialoger.ShowIt("Invalid code to recover")
		return false
	}

	// Create new password.
	newPassword, _ := r.Dialoger.GetSomeThing(
		fmt.Sprintf("%s",
			"Enter a new password..."))

	updateUser := &model.User{Password: newPassword, Email: user.GetModelUser().Email}
	info, err := r.AuthService.Recovery(user, updateUser)

	if err != nil {
		r.Dialoger.ShowIt(info)
		return false
	}

	r.Dialoger.ShowIt("Recovery password is successful")
	return true
}

// func (r *RootAuth) ExecuteAuth(authService auth.Auth, user user.User) bool {
// 	var info string

// 	// Identification.
// 	if ok := authService.Identification(user); ok {

// 		// Authentication. Verification
// 		email, password, err := r.Dialoger.VerifyDataAuth(user)

// 		if err == nil {
// 			checkUser := &model.User{Email: email, Password: password}
// 			info, err = authService.Authentication(user, checkUser)
// 			if err == nil {
// 				r.Dialoger.ShowIt("Authentication is successful")
// 				return true
// 			}
// 			r.Dialoger.ShowIt(info)
// 		}

// 		// Recovery Password.
// 		if err == errs.ErrPassword {
// 			comm, _ := r.Dialoger.GetSomeThing(
// 				fmt.Sprintf("%s\n\r%s",
// 					"Do you want to recover your password? \n\r\t yes \n\r\t no",
// 					"go out: exit, need help: help"))

// 			if comm == "yes" {
// 				// Отправляем на обработку в сервис восстановления пароля.
// 				r.MailChan <- email

// 				codeUser, _ := r.Dialoger.GetSomeThing(
// 					fmt.Sprintf("%s",
// 						"Enter code to recover from your mail..."))

// 				// CTX сюда с time
// 				select {
// 				case <-r.Ctx.Done():
// 					r.Dialoger.ShowIt("Some problems with data recovery")

// 				case codeGener := <-r.CodeChan:
// 					// Check user's code with generation code
// 					if codeUser != codeGener {
// 						r.Dialoger.ShowIt("Invalid code to recover")

// 					} else {
// 						newPassword, _ := r.Dialoger.GetSomeThing(
// 							fmt.Sprintf("%s",
// 								"Enter a new password..."))

// 						updateUser := &model.User{Password: newPassword, Email: email}
// 						info, err := authService.Recovery(user, updateUser)

// 						if err == nil {
// 							r.Dialoger.ShowIt("Recovery password is successful")
// 							return true
// 						}
// 						r.Dialoger.ShowIt(info)
// 						return false
// 					}
// 				}
// 			}
// 		}
// 	}

// 	// Registration.
// 	newUser := model.NewUser(r.Dialoger.GetDataRegister())
// 	info, err := authService.Registration(user, newUser)
// 	if err == nil {
// 		r.Dialoger.ShowIt("Registration is successful")
// 		return true
// 	}
// 	r.Dialoger.ShowIt(info)
// 	return false
// }

func (r *RootAuth) ExecuteAuth(workOK bool, user user.User) (bool, error) {
	if !workOK {
		return false, nil
	}

	// Verification.
	email, password, err := r.Dialoger.VerifyDataAuth(user)
	if err != nil {
		return false, err
	}

	checkUser := &model.User{Email: email, Password: password}
	info, err := r.AuthService.Authentication(user, checkUser)
	if err != nil {
		r.Dialoger.ShowIt(info)
		return false, err
	}

	r.Dialoger.ShowIt("Authentication is successful")
	return true, nil
}

func (r *RootAuth) ExecuteRegistr(user user.User) bool {
	newUser := model.NewUser(r.Dialoger.GetDataRegister())
	info, err := r.AuthService.Registration(user, newUser)

	if err == nil {
		r.Dialoger.ShowIt("Registration is successful")
		return true
	}

	r.Dialoger.ShowIt(info)
	return false
}

func (r *RootAuth) getAnswerFromRecoveryService(ctx context.Context) string {
	ctxT, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	select {
	case <-ctxT.Done():
		return ""
	case code := <-r.CodeChan:
		return code
	}
}
