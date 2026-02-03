package comm

import (
	"fmt"

	"github.com/boginskiy/GophKeeper/client/cmd/client"
	"github.com/boginskiy/GophKeeper/client/internal/cli"
	"github.com/boginskiy/GophKeeper/client/internal/user"
)

type CommImage struct {
	DialogSrv cli.ShowGetter
}

func NewCommImage(dialog cli.ShowGetter) *CommImage {
	return &CommImage{DialogSrv: dialog}
}

func (c *CommImage) Registration(client *client.ClientCLI, user *user.UserCLI) {
authLoop:
	for {
		// Get command.
		comm, _ := c.DialogSrv.GetSomeThing(

			fmt.Sprintf("%s\n\r%s",
				"What do you want to do with the image: \n\r\t create \n\r\t read \n\r\t update \n\r\t delete",
				"come back: back, need help: help"))

		switch comm {
		case "back", "help":
			break authLoop

		case "create":
			// r.RoutText.Execute(client, user)

		case "read":
			// r.RoutText.Execute(client, user)

		case "update":
			// r.RoutText.Execute(client, user)

		case "delete":
			// r.RoutText.Execute(client, user)

		default:
			c.DialogSrv.ShowIt("Invalid command. Try again...")
		}
	}
}
