package cli

import (
	"fmt"

	"github.com/boginskiy/GophKeeper/client/cmd/client"
	"github.com/boginskiy/GophKeeper/client/cmd/config"
	"github.com/boginskiy/GophKeeper/client/internal/errs"
	"github.com/boginskiy/GophKeeper/client/internal/logg"
	"github.com/boginskiy/GophKeeper/client/internal/user"
	"github.com/boginskiy/GophKeeper/client/pkg"
)

type GetterFn func(*client.ClientCLI, *user.UserCLI) (string, error)
type CheckerFn func(string, string) bool

type DialogService struct {
	Cfg  config.Config
	Logg logg.Logger
}

func NewDialogService(cfg config.Config, logger logg.Logger) *DialogService {
	return &DialogService{
		Cfg:  cfg,
		Logg: logger,
	}
}

// ShowHello
func (d *DialogService) ShowHello(client *client.ClientCLI, user *user.UserCLI) {
	var tx string
	if user.User != nil {
		tx = fmt.Sprintf("Hello, %s!", user.User.UserName)
	} else {
		tx = "Hello, Man!"
	}
	client.SendMess(tx)
}

func (d *DialogService) ShowIt(client *client.ClientCLI, it string) {
	client.SendMess(it)
}

func (d *DialogService) ShowErr(client *client.ClientCLI, err error) {
	client.SendMess(err.Error())
}

// GetSomeThing
func (d *DialogService) GetSomeThing(client *client.ClientCLI, user *user.UserCLI, mess string) (string, error) {
	client.SendMess(mess)
	return user.ReceiveMess()
}

// GetIt gives us everything we ask for.
func (d *DialogService) GetIt(client *client.ClientCLI, user *user.UserCLI, it string) (string, error) {
	client.SendMess(fmt.Sprintf("Enter the %s...", it))
	return user.ReceiveMess()
}

// GetEmail
func (d *DialogService) GetEmail(client *client.ClientCLI, user *user.UserCLI) (string, error) {
	client.SendMess("Enter the email...")
	return user.ReceiveMess()
}

// GetEmailWithCheck decorator GetEmail.
func (d *DialogService) GetEmailWithCheck(dialogGet GetterFn, funcCheck CheckerFn) GetterFn {
	return func(client *client.ClientCLI, user *user.UserCLI) (string, error) {

		for repeat := 0; repeat < d.Cfg.GetAttempts(); repeat++ {
			result, err := dialogGet(client, user)

			if err == nil && funcCheck(user.User.Email, result) {
				return result, nil
			}
			client.SendMess("Uncorrected credentials. Try again...")
		}
		return "", errs.ErrUncorrectCredentials
	}
}

// GetPassword
func (d *DialogService) GetPassword(client *client.ClientCLI, user *user.UserCLI) (string, error) {
	client.SendMess("Enter password...")
	return user.ReceiveMess()
}

// GetPasswordWithCheck decorator for GetPassword.
func (d *DialogService) GetPasswordWithCheck(dialogGet GetterFn, funcCheck CheckerFn) GetterFn {
	return func(client *client.ClientCLI, user *user.UserCLI) (string, error) {

		for repeat := 0; repeat < d.Cfg.GetAttempts(); repeat++ {
			result, err := dialogGet(client, user)

			if err == nil && funcCheck(user.User.Password, result) {
				return result, nil
			}
			client.SendMess("Uncorrected credentials. Try again...")
		}
		return "", errs.ErrUncorrectCredentials
	}
}

// DialogsAbRegister return many args.
func (d *DialogService) DialogsAbRegister(client *client.ClientCLI, user *user.UserCLI) (userName, email, phone, password string) {
	d.ShowIt(client, "You need to register...")

	userName, err := d.GetIt(client, user, "user name")
	d.Logg.CheckWithFatal(err, "bad user name")

	email, err = d.GetEmail(client, user)
	d.Logg.CheckWithFatal(err, "bad email")

	phone, err = d.GetIt(client, user, "phone")
	d.Logg.CheckWithFatal(err, "bad phone")

	password, err = d.GetPassword(client, user)
	d.Logg.CheckWithFatal(err, "bad password")

	return userName, email, phone, password
}

func (d *DialogService) DialogsAbAction(client *client.ClientCLI, user *user.UserCLI, action string) string {
	question := fmt.Sprintf(
		"%s to %s: %s",
		"What type of data do you want", action, "\n\r\t cred \n\r\t phone \n\r\t card \n\r\t info")

	hint := "come back: back, need help: help, pass: enter"

	result, _ := d.GetSomeThing(client, user, fmt.Sprintf("%s\n\r%s", question, hint))
	return result
}

// Checker
func (d *DialogService) CheckEmail(userEmail, email string) bool {
	return userEmail == email
}

func (d *DialogService) CheckPassword(userPassword, password string) bool {
	return pkg.CompareHashAndPassword(userPassword, password)
}
