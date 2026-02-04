package comm

import (
	"fmt"

	"github.com/boginskiy/GophKeeper/client/internal/errs"
	"github.com/boginskiy/GophKeeper/client/internal/infra"
	"github.com/boginskiy/GophKeeper/client/internal/rpc"
	"github.com/boginskiy/GophKeeper/client/internal/service"
	"github.com/boginskiy/GophKeeper/client/internal/user"
	"google.golang.org/grpc/metadata"
)

type CommMedia struct {
	Checker   infra.Checker
	DialogSrv infra.ShowGetter
	Service   service.BytesServicer
}

func NewCommMedia(
	checker infra.Checker,
	dialog infra.ShowGetter,
	srv service.BytesServicer) *CommMedia {

	return &CommMedia{
		Checker:   checker,
		DialogSrv: dialog,
		Service:   srv}
}

func (c *CommMedia) Registration(user user.User, dataType string) {
authLoop:
	for {
		// Get command.
		comm, _ := c.DialogSrv.GetSomeThing(

			fmt.Sprintf("%s\n\r%s",
				"What do you want to do with the media: \n\r\t upload \n\r\t unload \n\r\t read \n\r\t read-all \n\r\t delete",
				"come back: back, need help: help"))

		switch comm {
		case "back", "help":
			break authLoop

		case "upload":
			c.executeUpload(user, dataType)
		case "unload":
			c.executeUnload(user, dataType)
		case "read":
			c.executeRead(user, dataType)
		case "read-all":
			c.executeReadAll(user, dataType)
		case "delete":
			c.executeDelete(user, dataType)

		// case "update":
		// 	c.executeUpdate(user)

		default:
			c.DialogSrv.ShowIt("Invalid command. Try again...")
		}
	}
}

func (c *CommMedia) executeUpload(user user.User, dataType string) {
	// Call window.
	pathToFile, err := c.DialogSrv.CallWindows("Enter the abs path to file...")
	if err != nil {
		c.DialogSrv.ShowErr(err)
		return
	}

	// Check file.
	ok := c.Checker.CheckTypeFile(pathToFile, dataType)
	if !ok {
		c.DialogSrv.ShowIt("Invalid file type")
		return
	}

	// Call Service.
	obj, err := c.Service.Upload(user, pathToFile, dataType)
	if err != nil {
		c.DialogSrv.ShowErr(err)
		return
	}

	res, ok := obj.(*rpc.UploadBytesResponse)
	if !ok {
		c.DialogSrv.ShowErr(errs.ErrTypeConversion)
		return
	}

	c.DialogSrv.ShowIt(
		fmt.Sprintf("%s %s  sent size: %s received: %s\n\r",
			res.Status, res.UpdatedAt, res.SentSize, res.ReceivedSize))
}

func (c *CommMedia) executeUnload(user user.User, _ string) {
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

func (c *CommMedia) executeRead(user user.User, _ string) {
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
		fmt.Sprintf("%s   type: %s   created: %s\n\r", nameFile, res.Type, res.CreatedAt))
}

func (c *CommMedia) executeReadAll(user user.User, dataType string) {
	// Call Service.
	obj, err := c.Service.ReadAll(user, dataType)
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

func (c *CommMedia) executeDelete(user user.User, _ string) {
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
