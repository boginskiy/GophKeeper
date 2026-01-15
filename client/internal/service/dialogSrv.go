package service

import (
	"fmt"

	"github.com/boginskiy/GophKeeper/client/internal/client"
	"github.com/boginskiy/GophKeeper/client/internal/user"
)

type DialogService struct {
	Client *client.ClientCLI
	User   *user.UserCLI
}

func NewDialogService(clt *client.ClientCLI, usr *user.UserCLI) *DialogService {
	return &DialogService{
		Client: clt,
		User:   usr,
	}
}

func (d *DialogService) Hello() {
	d.Client.SendMess("Hello! Press the 'Enter'...")
	d.User.ReceiveMess()
}

func (d *DialogService) Registration() {
	d.Client.SendMess("You need to register...")
	d.Client.SendMess("You need to register...")

	// Какие данные ?
	// 
}

func (d *DialogService) Run() {
	d.Hello()
	d.Registration()

	// Проверка регистрации пользователя, если данные подгрузились.
	// Если нет данных, то нужно пользователя зарегать, как локально, так и на удаленном сервере

	fmt.Println("end <<")
}
