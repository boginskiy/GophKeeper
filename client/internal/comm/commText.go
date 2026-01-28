package comm

import (
	"fmt"

	"github.com/boginskiy/GophKeeper/client/cmd/client"
	"github.com/boginskiy/GophKeeper/client/internal/cli"
	"github.com/boginskiy/GophKeeper/client/internal/service"
	"github.com/boginskiy/GophKeeper/client/internal/user"
)

type CommText struct {
	Dialoger cli.Dialoger
	Service  service.Servicer
}

func NewCommText(dialoger cli.Dialoger, srv service.Servicer) *CommText {
	return &CommText{Dialoger: dialoger, Service: srv}
}

func (c *CommText) Registration(client *client.ClientCLI, user *user.UserCLI) {
authLoop:
	for {
		// Get command.
		comm, _ := c.Dialoger.GetSomeThing(client, user,
			fmt.Sprintf("%s\n\r%s",
				"What do you want to do with the text: \n\r\t create \n\r\t read \n\r\t read-all \n\r\t update \n\r\t delete",
				"come back: back, need help: help"))

		switch comm {
		case "back", "help":
			break authLoop

		case "create":
			c.Service.Create(client, user)
		case "read":
			c.Service.Read(client, user)
		case "read-all":
			c.Service.ReadAll(client, user)
		case "update":
			c.Service.Update(client, user)
		case "delete":
			c.Service.Delete(client, user)

		default:
			c.Dialoger.ShowIt(client, "Invalid command. Try again...")
		}
	}
}
