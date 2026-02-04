package comm

import (
	"fmt"

	"github.com/boginskiy/GophKeeper/client/internal/errs"
	"github.com/boginskiy/GophKeeper/client/internal/infra"
	"github.com/boginskiy/GophKeeper/client/internal/rpc"
	"github.com/boginskiy/GophKeeper/client/internal/service"
	"github.com/boginskiy/GophKeeper/client/internal/user"
	"github.com/boginskiy/GophKeeper/client/pkg"
)

type CommVideo struct {
	Checker   infra.Checker
	DialogSrv infra.ShowGetter
	Service   service.BytesServicer
	Tp        string
}

func NewCommVideo(checker infra.Checker, dialog infra.ShowGetter, srv service.BytesServicer) *CommVideo {
	return &CommVideo{Checker: checker, DialogSrv: dialog, Service: srv, Tp: "video"}
}

func (c *CommVideo) Registration(user user.User, dataType string) {
authLoop:
	for {
		// Get command.
		comm, _ := c.DialogSrv.GetSomeThing(

			fmt.Sprintf("%s\n\r%s",
				"What do you want to do with the video: \n\r\t upload \n\r\t unload \n\r\t read-all \n\r\t delete",
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

func (c *CommVideo) executeUpload(user user.User) {
	// Вызываем окно для получения пути к загружаемому файлу
	pathToFile, err := pkg.Selector("selector.py")
	if err != nil || pathToFile == "" {
		// Если ошибка вызова окна, то вводим вручную путь до файла.
		pathToFile, _ = c.DialogSrv.GetSomeThing("Enter the abs path to video file...")
	}

	// Проверка что выбранный файл типа sound.
	ok := c.Checker.CheckTypeFile(pathToFile, c.Tp)
	if !ok {
		c.DialogSrv.ShowIt("There must be a video file!")
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

func (c *CommVideo) executeUnload(user user.User) {

}

func (c *CommVideo) executeReadAll(user user.User) {

}

func (c *CommVideo) executeDelete(user user.User) {

}
