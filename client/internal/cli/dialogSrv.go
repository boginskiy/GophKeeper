package cli

import (
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

// ShowSomeInfo
func (d *DialogService) ShowSomeInfo(client *client.ClientCLI, info string) {
	client.SendMess(info)
}

// ShowHello
func (d *DialogService) ShowHello(client *client.ClientCLI, user *user.UserCLI) {
	client.SendMess("Hello! Press the 'Enter'...")
	user.ReceiveMess()
}

// ShowLogIn
func (d *DialogService) ShowLogIn(client *client.ClientCLI, user *user.UserCLI) {
	client.SendMess("You need log in...")
}

// ShowRegister
func (d *DialogService) ShowRegister(client *client.ClientCLI, user *user.UserCLI) {
	client.SendMess("You need to register...")
}

// ShowStatusAuth
func (d *DialogService) ShowStatusAuth(client *client.ClientCLI, user *user.UserCLI) {
	client.SendMess("Authentication is successful")
}

// GetUserName
func (d *DialogService) GetUserName(client *client.ClientCLI, user *user.UserCLI) (string, error) {
	client.SendMess("Enter the user name...")
	return user.ReceiveMess()
}

// GetUserName
func (d *DialogService) GetSomeThing(client *client.ClientCLI, user *user.UserCLI, mess string) (string, error) {
	client.SendMess(mess)
	return user.ReceiveMess()
}

// GetPhone
func (d *DialogService) GetPhone(client *client.ClientCLI, user *user.UserCLI) (string, error) {
	client.SendMess("Enter the phone...")
	return user.ReceiveMess()
}

// ShowHello
func (d *DialogService) GetCommand(client *client.ClientCLI, user *user.UserCLI) (string, error) {
	client.SendMess("Enter the command...")
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

// Checker
func (d *DialogService) CheckEmail(userEmail, email string) bool {
	return userEmail == email
}

func (d *DialogService) CheckPassword(userPassword, password string) bool {
	return pkg.CompareHashAndPassword(userPassword, password)
}

func (d *DialogService) DialogsAbRegister(client *client.ClientCLI, user *user.UserCLI) (userName, email, phone, password string) {
	d.ShowRegister(client, user)

	userName, err := d.GetUserName(client, user)
	d.Logg.CheckWithFatal(err, "bad user name")

	email, err = d.GetEmail(client, user)
	d.Logg.CheckWithFatal(err, "bad email")

	phone, err = d.GetPhone(client, user)
	d.Logg.CheckWithFatal(err, "bad phone")

	password, err = d.GetPassword(client, user)
	d.Logg.CheckWithFatal(err, "bad password")

	return userName, email, phone, password
}
