package comm

import (
	"fmt"

	"github.com/boginskiy/GophKeeper/client/cmd/client"
	"github.com/boginskiy/GophKeeper/client/internal/cli"
	"github.com/boginskiy/GophKeeper/client/internal/user"
)

type CommSound struct {
	Dialoger cli.Dialoger
}

func NewCommSound(dialoger cli.Dialoger) *CommSound {
	return &CommSound{Dialoger: dialoger}
}

func (c *CommSound) Registration(client *client.ClientCLI, user *user.UserCLI) {
authLoop:
	for {
		// Get command.
		comm, _ := c.Dialoger.GetSomeThing(client, user,
			fmt.Sprintf("%s\n\r%s",
				"What do you want to do with the sound: \n\r\t create \n\r\t read \n\r\t update \n\r\t delete",
				"go to previous step: .., need help: help"))

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
			c.Dialoger.ShowSomeInfo(client, "Invalid command. Try again...")
		}
	}
}
