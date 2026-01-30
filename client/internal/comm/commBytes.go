package comm

import (
	"fmt"

	"github.com/boginskiy/GophKeeper/client/cmd/client"
	"github.com/boginskiy/GophKeeper/client/internal/cli"
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
			c.Service.Upload(client, user)
		case "unload":
			c.Service.Unload(client, user)
		// case "update":
		// 	c.Service.Update(client, user)
		// case "delete":
		// 	c.Service.Delete(client, user)

		default:
			c.DialogSrv.ShowIt("Invalid command. Try again...")
		}
	}
}
