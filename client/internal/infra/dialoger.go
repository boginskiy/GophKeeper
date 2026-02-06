package infra

import (
	"fmt"

	"github.com/boginskiy/GophKeeper/client/cmd/client"
	"github.com/boginskiy/GophKeeper/client/cmd/config"
	"github.com/boginskiy/GophKeeper/client/internal/errs"
	"github.com/boginskiy/GophKeeper/client/internal/logg"
	"github.com/boginskiy/GophKeeper/client/internal/user"
	"github.com/boginskiy/GophKeeper/client/pkg"
)

type Dialog struct {
	Cfg     config.Config
	Logg    logg.Logger
	Checker Checker
	Client  client.Client
	User    user.User
}

func NewDialog(
	cfg config.Config,
	logger logg.Logger,
	ch Checker,
	cl client.Client,
	us user.User,
) *Dialog {

	return &Dialog{
		Cfg:     cfg,
		Logg:    logger,
		Checker: ch,
		Client:  cl,
		User:    us,
	}
}

func (d *Dialog) ShowIt(it string) {
	d.Client.SendMess(it)
}

func (d *Dialog) ShowErr(err error) {
	d.Client.SendMess(err.Error())
}

// GetSomeThing
func (d *Dialog) GetSomeThing(mess string) (string, error) {
	d.Client.SendMess(mess)
	return d.User.ReceiveMess()
}

// GetEnterIt gives us everything we ask for.
func (d *Dialog) GetEnterIt(it string) (string, error) {
	d.Client.SendMess(fmt.Sprintf("Enter the %s...", it))
	return d.User.ReceiveMess()
}

func (d *Dialog) GetDataAction(action string) string {
	question := fmt.Sprintf(
		"%s to %s %s",
		"What type of text data do you want", action, "\n\r\t info \n\r\t phone \n\r\t card \n\r\t other")

	hint := "come back: back, need help: help, pass: enter"
	result, _ := d.GetSomeThing(fmt.Sprintf("%s\n\r%s", question, hint))
	return result
}

func (d *Dialog) VerifyEnterIt(needToTake, needToCompare string, quantity int) (string, error) {
	for q := 0; q < quantity; q++ {
		result, err := d.GetEnterIt(needToTake)
		if err == nil && d.Checker.CheckTwoString(needToCompare, result) {
			return result, nil
		}
		d.Client.SendMess("Uncorrected credentials. Try again...")
	}
	return "", errs.ErrEmail
}

func (d *Dialog) VerifyEnterPassword(needToTake, needToCompare string, quantity int) (string, error) {
	for q := 0; q < quantity; q++ {
		result, err := d.GetEnterIt(needToTake)

		if err == nil && d.Checker.CheckPassword(needToCompare, result) {
			return result, nil
		}
		d.Client.SendMess("Uncorrected credentials. Try again...")
	}
	return "", errs.ErrPassword
}

func (d *Dialog) GetDataRegister() (userName, email, phone, password string) {
	d.ShowIt("You need to register...")

	userName, err := d.GetEnterIt("user name")
	d.Logg.CheckWithFatal(err, "bad user name")

	email, err = d.GetEnterIt("email")
	d.Logg.CheckWithFatal(err, "bad email")

	phone, err = d.GetEnterIt("phone")
	d.Logg.CheckWithFatal(err, "bad phone")

	password, err = d.GetEnterIt("password")
	d.Logg.CheckWithFatal(err, "bad password")

	return userName, email, phone, password
}

func (d *Dialog) VerifyDataAuth(user user.User) (email, password string, err error) {
	d.ShowIt("You need log in")

	email, err = d.VerifyEnterIt("email", user.GetModelUser().Email, d.Cfg.GetMaxRetries())
	if err != nil {
		return "", "", err
	}

	password, err = d.VerifyEnterPassword("password", user.GetModelUser().Password, d.Cfg.GetMaxRetries())
	if err != nil {
		return "", "", err
	}

	return email, password, nil
}

func (d *Dialog) CallWindows(input string) (string, error) {
	// Вызываем окно для получения пути к загружаемому файлу
	pathToFile, err := pkg.Selector("selector.py")

	// Если ошибка вызова окна, то вводим вручную путь до файла.
	if err != nil || pathToFile == "" {
		pathToFile, err = d.GetSomeThing(input)
	}

	if err != nil {
		return "", err
	}
	return pathToFile, nil
}
