package service

import (
	"fmt"

	"github.com/boginskiy/GophKeeper/client/internal/client"
	"github.com/boginskiy/GophKeeper/client/internal/logg"
	"github.com/boginskiy/GophKeeper/client/internal/user"
)

type DialogService struct {
	Logger logg.Logger
	Client *client.ClientCLI
	User   *user.UserCLI
}

func NewDialogService(logger logg.Logger, clt *client.ClientCLI, usr *user.UserCLI) *DialogService {
	return &DialogService{
		Logger: logger,
		Client: clt,
		User:   usr,
	}
}

func (d *DialogService) Welcome() {
	d.Client.SendMess("Hello! Press the 'Enter'...")
	d.User.ReceiveMess()
}

func (d *DialogService) Authorization() {
	// email
	d.Client.SendMess(
		"You need log in...",
		"Enter email...")
	email, err := d.User.ReceiveMess()
	d.Logger.CheckWithFatal(err, "bad email")

	// password
	d.Client.SendMess("Enter password...")
	password, err := d.User.ReceiveMess()
	d.Logger.CheckWithFatal(err, "bad password")

	// TODO
	// Нужна проверка введенного пароля

	fmt.Println(email, password)
}

func (d *DialogService) Registration() {
	// user name
	d.Client.SendMess(
		"You need to register...",
		"Enter user name...")
	userName, err := d.User.ReceiveMess()
	// TODO нужна проверка
	d.Logger.CheckWithFatal(err, "bad user name")

	// email
	d.Client.SendMess("Enter email...")
	email, err := d.User.ReceiveMess()
	// TODO нужна проверка c запросом на сервер
	d.Logger.CheckWithFatal(err, "bad email")

	// phone
	d.Client.SendMess("Enter phone...")
	phone, err := d.User.ReceiveMess()
	// TODO нужна проверка
	d.Logger.CheckWithFatal(err, "bad phone")

	// password
	d.Client.SendMess("Enter password...")
	password, err := d.User.ReceiveMess()
	// TODO нужна проверка
	d.Logger.CheckWithFatal(err, "bad password")

	d.User.NewUser(userName, email, phone, password)
	d.Client.SendMess("Registration is successful")
}

func (d *DialogService) Run() {
	// Приветствие.
	d.Welcome()

	// Регистрация/Авторизация пользователя.
	if d.User.User != nil {
		d.Authorization()
	} else {
		d.Registration()
	}

	// Конец CLI сессии.
	fmt.Println(">> end <<")
}

// TODO
// Если есть данные н апользователя, просим ввести почту и пароль.
// Пароль надо хешировать
