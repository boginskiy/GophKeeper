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
	Cfg    config.Config
	Logger logg.Logger
}

func NewDialogService(cfg config.Config, logger logg.Logger) *DialogService {
	return &DialogService{
		Cfg:    cfg,
		Logger: logger,
	}
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

// GetEmail
func (d *DialogService) GetEmail(client *client.ClientCLI, user *user.UserCLI) (string, error) {
	client.SendMess("Enter email...")
	return user.ReceiveMess()
}

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

// GetUserName
func (d *DialogService) GetUserName(client *client.ClientCLI, user *user.UserCLI) (string, error) {
	client.SendMess("Enter user name...")
	return user.ReceiveMess()
}

// GetPhone
func (d *DialogService) GetPhone(client *client.ClientCLI, user *user.UserCLI) (string, error) {
	client.SendMess("Enter phone...")
	return user.ReceiveMess()
}

// Checker
func (d *DialogService) CheckEmail(userEmail, email string) bool {
	return userEmail == email
}

func (d *DialogService) CheckPassword(userPassword, password string) bool {
	return pkg.CompareHashAndPassword(userPassword, password)
}
