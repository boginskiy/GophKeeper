package cli

import (
	"fmt"

	"github.com/boginskiy/GophKeeper/client/cmd/client"
	"github.com/boginskiy/GophKeeper/client/cmd/config"
	"github.com/boginskiy/GophKeeper/client/internal/logg"
	"github.com/boginskiy/GophKeeper/client/internal/user"
	"github.com/boginskiy/GophKeeper/client/internal/utils"
)

type DialogService struct {
	Cfg     config.Config
	Logg    logg.Logger
	Checker utils.Checker
	Client  client.Client
	User    user.User
}

func NewDialogService(
	cfg config.Config,
	logger logg.Logger,
	ch utils.Checker,
	cl client.Client,
	us user.User,
) *DialogService {

	return &DialogService{
		Cfg:     cfg,
		Logg:    logger,
		Checker: ch,
		Client:  cl,
		User:    us,
	}
}

func (d *DialogService) ShowIt(it string) {
	d.Client.SendMess(it)
}

func (d *DialogService) ShowErr(err error) {
	d.Client.SendMess(err.Error())
}

// GetSomeThing
func (d *DialogService) GetSomeThing(mess string) (string, error) {
	d.Client.SendMess(mess)
	return d.User.ReceiveMess()
}

// GetEnterIt gives us everything we ask for.
func (d *DialogService) GetEnterIt(it string) (string, error) {
	d.Client.SendMess(fmt.Sprintf("Enter the %s...", it))
	return d.User.ReceiveMess()
}

func (d *DialogService) GetDataAction(action string) string {
	question := fmt.Sprintf(
		"%s to %s %s",
		"What type of text data do you want", action, "\n\r\t info \n\r\t phone \n\r\t card \n\r\t other")

	hint := "come back: back, need help: help, pass: enter"
	result, _ := d.GetSomeThing(fmt.Sprintf("%s\n\r%s", question, hint))
	return result
}

func (d *DialogService) VerifyEnterIt(needToTake, needToCompare string, quantity int) bool {
	for q := 0; q < quantity; q++ {
		result, err := d.GetEnterIt(needToTake)
		if err == nil && d.Checker.CheckTwoString(needToCompare, result) {
			return true
		}
		d.Client.SendMess("Uncorrected credentials. Try again...")
	}
	return false
}

func (d *DialogService) VerifyEnterPassword(needToTake, needToCompare string, quantity int) bool {
	for q := 0; q < quantity; q++ {
		result, err := d.GetEnterIt(needToTake)

		if err == nil && d.Checker.CheckPassword(needToCompare, result) {
			return true
		}
		d.Client.SendMess("Uncorrected credentials. Try again...")
	}
	return false
}

func (d *DialogService) GetDataRegister() (userName, email, phone, password string) {
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

func (d *DialogService) VerifyDataAuth(user user.User) bool {
	d.ShowIt("You need log in")
	ok := d.VerifyEnterIt("email", user.GetModelUser().Email, d.Cfg.GetAttempts())
	ok2 := d.VerifyEnterPassword("password", user.GetModelUser().Password, d.Cfg.GetAttempts())

	return ok && ok2
}
