package comm

import (
	"fmt"

	"github.com/boginskiy/GophKeeper/client/cmd/client"
	"github.com/boginskiy/GophKeeper/client/internal/cli"
	"github.com/boginskiy/GophKeeper/client/internal/infra"
	"github.com/boginskiy/GophKeeper/client/internal/service"
	"github.com/boginskiy/GophKeeper/client/internal/user"
)

type CommVideo struct {
	Checker   infra.Checker
	DialogSrv cli.ShowGetter
	Service   service.BytesServicer
	Tp        string
}

func NewCommVideo(checker infra.Checker, dialog cli.ShowGetter, srv service.BytesServicer) *CommVideo {
	return &CommVideo{Checker: checker, DialogSrv: dialog, Service: srv, Tp: "video"}
}

func (c *CommVideo) Registration(client *client.ClientCLI, user *user.UserCLI) {
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
