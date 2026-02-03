package comm

import (
	"fmt"

	"github.com/boginskiy/GophKeeper/client/cmd/client"
	"github.com/boginskiy/GophKeeper/client/internal/cli"
	"github.com/boginskiy/GophKeeper/client/internal/errs"
	"github.com/boginskiy/GophKeeper/client/internal/infra"
	"github.com/boginskiy/GophKeeper/client/internal/rpc"
	"github.com/boginskiy/GophKeeper/client/internal/service"
	"github.com/boginskiy/GophKeeper/client/internal/user"
	"google.golang.org/grpc/metadata"
)

type CommBytes struct {
	DialogSrv cli.ShowGetter
	Service   service.BytesServicer
	Tp        string
}

func NewCommBytes(dialog cli.ShowGetter, srv service.BytesServicer) *CommBytes {
	return &CommBytes{DialogSrv: dialog, Service: srv, Tp: "bytes"}
}

func (c *CommBytes) Registration(client *client.ClientCLI, user *user.UserCLI) {
authLoop:
	for {
		// Get command.
		comm, _ := c.DialogSrv.GetSomeThing(

			fmt.Sprintf("%s\n\r%s",
				"What do you want to do with the bytes: \n\r\t upload \n\r\t unload \n\r\t read \n\r\t read-all \n\r\t delete",
				"come back: back, need help: help"))

		switch comm {
		case "back", "help":
			break authLoop

		case "upload":
			c.executeUpload(user)
		case "unload":
			c.executeUnload(user)
		case "read":
			c.executeRead(user)
		case "read-all":
			c.executeReadAll(user)
		case "delete":
			c.executeDelete(user)

		// case "update":
		// 	c.executeUpdate(user)

		default:
			c.DialogSrv.ShowIt("Invalid command. Try again...")
		}
	}
}

func (c *CommBytes) executeUpload(user user.User) {
	pathToFile, _ := c.DialogSrv.GetSomeThing("Enter the abs path to file...")

	// Call Service.
	obj, err := c.Service.Upload(user, pathToFile, c.Tp)

	if err != nil {
		c.DialogSrv.ShowErr(err)
		return
	}
	// TODO...
	// Ну вот говно же полное, что ты не возвращаешь универсальные данные
	// у каждого метода Service свой тип ответа.
	res, ok := obj.(*rpc.UploadBytesResponse)
	if !ok {
		c.DialogSrv.ShowErr(errs.ErrTypeConversion)
		return
	}

	c.DialogSrv.ShowIt(
		fmt.Sprintf("%s %s  sent size: %s received: %s\n\r",
			res.Status, res.UpdatedAt, res.SentSize, res.ReceivedSize))
}

func (c *CommBytes) executeUnload(user user.User) {
	nameFile, _ := c.DialogSrv.GetSomeThing("Enter the name file...")

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

func (c *CommBytes) executeRead(user user.User) {
	nameFile, _ := c.DialogSrv.GetSomeThing("Enter the name file...")

	// Call Service.
	obj, err := c.Service.Read(user, nameFile)

	if err != nil {
		c.DialogSrv.ShowErr(err)
		return
	}

	res, ok := obj.(*rpc.ReadBytesResponse)
	if !ok {
		c.DialogSrv.ShowErr(errs.ErrTypeConversion)
		return
	}

	c.DialogSrv.ShowIt(
		fmt.Sprintf("%s   type: %s   created: %s\n\r",
			nameFile, res.Type, res.CreatedAt))
}

func (c *CommBytes) executeReadAll(user user.User) {
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

func (c *CommBytes) executeDelete(user user.User) {
	nameFile, _ := c.DialogSrv.GetSomeThing("Enter the name file...")

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
