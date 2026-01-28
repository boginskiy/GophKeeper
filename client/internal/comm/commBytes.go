package comm

import (
	"fmt"

	"github.com/boginskiy/GophKeeper/client/cmd/client"
	"github.com/boginskiy/GophKeeper/client/internal/cli"
	"github.com/boginskiy/GophKeeper/client/internal/service"
	"github.com/boginskiy/GophKeeper/client/internal/user"
)

type CommBytes struct {
	Dialoger cli.Dialoger
	Service  service.ServicerByter
}

func NewCommBytes(dialoger cli.Dialoger, srv service.ServicerByter) *CommBytes {
	return &CommBytes{Dialoger: dialoger, Service: srv}
}

func (c *CommBytes) Registration(client *client.ClientCLI, user *user.UserCLI) {
authLoop:
	for {
		// Get command.
		comm, _ := c.Dialoger.GetSomeThing(client, user,
			fmt.Sprintf("%s\n\r%s",
				"What do you want to do with the bytes: \n\r\t create \n\r\t read \n\r\t update \n\r\t delete",
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
			c.Dialoger.ShowIt(client, "Invalid command. Try again...")
		}
	}
}
