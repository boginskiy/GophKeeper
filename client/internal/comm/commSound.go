package comm

import (
	"fmt"

	"github.com/boginskiy/GophKeeper/client/internal/errs"
	"github.com/boginskiy/GophKeeper/client/internal/infra"
	"github.com/boginskiy/GophKeeper/client/internal/rpc"
	"github.com/boginskiy/GophKeeper/client/internal/service"
	"github.com/boginskiy/GophKeeper/client/internal/user"
	"github.com/boginskiy/GophKeeper/client/pkg"
	"google.golang.org/grpc/metadata"
)

type CommSound struct {
	Checker   infra.Checker
	DialogSrv infra.ShowGetter
	Service   service.BytesServicer
	Tp        string
}

func NewCommSound(checker infra.Checker, dialog infra.ShowGetter, srv service.BytesServicer) *CommSound {
	return &CommSound{Checker: checker, DialogSrv: dialog, Service: srv, Tp: "sound"}
}

func (c *CommSound) Registration(user user.User, dataType string) {
authLoop:
	for {
		// Get command.
		comm, _ := c.DialogSrv.GetSomeThing(

			fmt.Sprintf("%s\n\r%s",
				"What do you want to do with the sound: \n\r\t upload \n\r\t unload \n\r\t read-all \n\r\t delete",
				"come back: back, need help: help"))

		switch comm {
		case "back", "help":
			break authLoop

		case "upload":
			c.executeUpload(user)
		case "unload":
			c.executeUnload(user)
		case "read-all":
			c.executeReadAll(user)
		case "delete":
			c.executeDelete(user)

		default:
			c.DialogSrv.ShowIt("Invalid command. Try again...")
		}
	}
}

func (c *CommSound) executeUpload(user user.User) {
	// Вызываем окно для получения пути к загружаемому файлу
	pathToFile, err := pkg.Selector("selector.py")
	if err != nil || pathToFile == "" {
		// Если ошибка вызова окна, то вводим вручную путь до файла.
		pathToFile, _ = c.DialogSrv.GetSomeThing("Enter the abs path to music file...")
	}

	// Проверка что выбранный файл типа sound.
	ok := c.Checker.CheckTypeFile(pathToFile, c.Tp)
	if !ok {
		c.DialogSrv.ShowIt("There must be a music file!")
		return
	}

	// Call Service.
	obj, err := c.Service.Upload(user, pathToFile, c.Tp)
	if err != nil {
		c.DialogSrv.ShowErr(err)
		return
	}

	res, ok := obj.(*rpc.UploadBytesResponse)
	if !ok {
		c.DialogSrv.ShowErr(errs.ErrTypeConversion)
		return
	}

	c.DialogSrv.ShowIt(fmt.Sprintf("%s %s  sent size: %s received: %s\n\r",
		res.Status, res.UpdatedAt, res.SentSize, res.ReceivedSize))
}

func (c *CommSound) executeUnload(user user.User) {
	nameFile, _ := c.DialogSrv.GetSomeThing("Enter the name music file...")

	// Проверка что выбранный файл типа sound.
	ok := c.Checker.CheckTypeFile(nameFile, c.Tp)
	if !ok {
		c.DialogSrv.ShowIt("There must be a music file!")
		return
	}

	// Call Service.
	obj, err := c.Service.Unload(user, nameFile)
	if err != nil {
		c.DialogSrv.ShowErr(err)
		return
	}

	serverHeader, ok := obj.(*metadata.MD)
	if !ok {
		c.DialogSrv.ShowErr(errs.ErrTypeConversion)
		return
	}

	c.DialogSrv.ShowIt(
		fmt.Sprintf("%s %s  sent size: %s received: %s\n\r",
			"unloaded",
			infra.TakeValueFromHeader(*serverHeader, "updated_at", 0),
			infra.TakeValueFromHeader(*serverHeader, "sent_size", 0),
			infra.TakeValueFromHeader(*serverHeader, "received_size", 0)))
}

func (c *CommSound) executeReadAll(user user.User) {
	// Call Service.
	obj, err := c.Service.ReadAll(user, c.Tp)
	if err != nil {
		c.DialogSrv.ShowErr(err)
		return
	}

	res, ok := obj.(*rpc.ReadAllBytesResponse)
	if !ok {
		c.DialogSrv.ShowErr(errs.ErrTypeConversion)
		return
	}

	for _, bytes := range res.BytesResponses {
		c.DialogSrv.ShowIt(fmt.Sprintf("%s   total size: %s   last update: %s", bytes.Name, bytes.TotalSize, bytes.UpdatedAt))
	}
	fmt.Println()
}

func (c *CommSound) executeDelete(user user.User) {
	nameFile, _ := c.DialogSrv.GetSomeThing("Enter the name music file...")

	// Проверка что выбранный файл типа sound.
	ok := c.Checker.CheckTypeFile(nameFile, c.Tp)
	if !ok {
		c.DialogSrv.ShowIt("There must be a music file!")
		return
	}

	// Call Service.
	obj, err := c.Service.Delete(user, nameFile)
	if err != nil {
		c.DialogSrv.ShowErr(err)
		return
	}

	res, ok := obj.(*rpc.DeleteBytesResponse)
	if !ok {
		c.DialogSrv.ShowErr(errs.ErrTypeConversion)
		return
	}

	c.DialogSrv.ShowIt(fmt.Sprintf("%s %s\n\r", res.Status, nameFile))
}
