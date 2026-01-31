package comm

import (
	"fmt"

	"github.com/boginskiy/GophKeeper/client/cmd/client"
	"github.com/boginskiy/GophKeeper/client/internal/cli"
	"github.com/boginskiy/GophKeeper/client/internal/errs"
	"github.com/boginskiy/GophKeeper/client/internal/rpc"
	"github.com/boginskiy/GophKeeper/client/internal/service"
	"github.com/boginskiy/GophKeeper/client/internal/user"
)

type CommBytes struct {
	DialogSrv cli.ShowGetter
	Service   service.ServicerByter
	Tp        string
}

func NewCommBytes(dialog cli.ShowGetter, srv service.ServicerByter) *CommBytes {
	return &CommBytes{DialogSrv: dialog, Service: srv, Tp: "bytes"}
}

func (c *CommBytes) Registration(client *client.ClientCLI, user *user.UserCLI) {
authLoop:
	for {
		// Get command.
		comm, _ := c.DialogSrv.GetSomeThing(

			fmt.Sprintf("%s\n\r%s",
				"What do you want to do with the bytes: \n\r\t upload \n\r\t unload \n\r\t update \n\r\t delete",
				"come back: back, need help: help"))

		switch comm {
		case "back", "help":
			break authLoop

		case "upload":
			c.executeUpload(user)
		case "unload":
			c.executeUnload(user)

		// case "update":
		// 	c.Service.Update(client, user)
		// case "delete":
		// 	c.Service.Delete(client, user)

		default:
			c.DialogSrv.ShowIt("Invalid command. Try again...")
		}
	}
}

func (c *CommBytes) executeUpload(user user.User) {
	pathToFile, _ := c.DialogSrv.GetSomeThing("Enter the abs path to file...")

	// Call Service.
	obj, err := c.Service.Upload(user, pathToFile)

	if err != nil {
		c.DialogSrv.ShowErr(err)
		return
	}

	res, ok := obj.(*rpc.UploadBytesResponse)
	if !ok {
		c.DialogSrv.ShowErr(errs.ErrTypeConversion)
		return
	}

	c.DialogSrv.ShowIt(fmt.Sprintf("%s %s\n\r", res.Status, res.UpdatedAt))
}

func (c *CommBytes) executeUnload(user user.User) {

}
