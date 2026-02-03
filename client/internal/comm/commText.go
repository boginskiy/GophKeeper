package comm

import (
	"fmt"

	"github.com/boginskiy/GophKeeper/client/cmd/client"
	"github.com/boginskiy/GophKeeper/client/internal/cli"
	"github.com/boginskiy/GophKeeper/client/internal/errs"
	"github.com/boginskiy/GophKeeper/client/internal/model"
	"github.com/boginskiy/GophKeeper/client/internal/rpc"
	"github.com/boginskiy/GophKeeper/client/internal/service"
	"github.com/boginskiy/GophKeeper/client/internal/user"
)

type CommText struct {
	DialogSrv cli.ShowGetter
	Service   service.TextServicer[model.Text]
	Tp        string
}

func NewCommText(
	dialog cli.ShowGetter,
	srv service.TextServicer[model.Text]) *CommText {

	return &CommText{
		DialogSrv: dialog,
		Service:   srv,
		Tp:        "text"}
}

func (c *CommText) Registration(client *client.ClientCLI, user *user.UserCLI) {
authLoop:
	for {
		// Get command.
		comm, _ := c.DialogSrv.GetSomeThing(

			fmt.Sprintf("%s\n\r%s",
				"What do you want to do with the text: \n\r\t create \n\r\t read \n\r\t read-all \n\r\t update \n\r\t delete",
				"come back: back, need help: help"))

		switch comm {
		case "back", "help":
			break authLoop

		case "create":
			c.executeCreate(user)
		case "read":
			c.executeRead(user)
		case "read-all":
			c.executeReadAll(user)
		case "update":
			c.executeUpdate(user)
		case "delete":
			c.executeDelete(user)

		default:
			c.DialogSrv.ShowIt("Invalid command. Try again...")
		}
	}
}

func (c *CommText) executeCreate(user user.User) {
	name := c.DialogSrv.GetDataAction("create")            // Get info about name text.
	tx, _ := c.DialogSrv.GetSomeThing("Enter the text...") // Enter text for saving.

	// Call Service.
	obj, err := c.Service.Create(user, *model.NewText(name, c.Tp, tx, user.GetModelUser().Email))

	if err != nil {
		c.DialogSrv.ShowErr(err)
		return
	}

	res, ok := obj.(*rpc.CreateResponse)
	if !ok {
		c.DialogSrv.ShowErr(errs.ErrTypeConversion)
		return
	}

	c.DialogSrv.ShowIt(fmt.Sprintf("%s %s\n\r", res.Status, res.UpdatedAt))
}

func (c *CommText) executeRead(user user.User) {
	name := c.DialogSrv.GetDataAction("read")

	// Call Service.
	obj, err := c.Service.Read(user, model.Text{Name: name, Owner: user.GetModelUser().Email})

	if err != nil {
		c.DialogSrv.ShowErr(err)
		return
	}

	res, ok := obj.(*rpc.ReadResponse)
	if !ok {
		c.DialogSrv.ShowErr(errs.ErrTypeConversion)
		return
	}

	c.DialogSrv.ShowIt(fmt.Sprintf("%s: %s   last update: %s\n\r", res.Name, res.Text, res.UpdatedAt))
}

func (c *CommText) executeReadAll(user user.User) {
	// Call Service.
	obj, err := c.Service.ReadAll(user, model.Text{Type: c.Tp, Owner: user.GetModelUser().Email})

	if err != nil {
		c.DialogSrv.ShowErr(err)
		return
	}

	res, ok := obj.(*rpc.ReadAllResponse)
	if !ok {
		c.DialogSrv.ShowErr(errs.ErrTypeConversion)
		return
	}

	for _, text := range res.TextResponses {
		c.DialogSrv.ShowIt(fmt.Sprintf("%s: %s   last update: %s", text.Name, text.Text, text.UpdatedAt))
	}
	fmt.Println()
}

func (c *CommText) executeUpdate(user user.User) {
	name := c.DialogSrv.GetDataAction("update")
	tx, _ := c.DialogSrv.GetSomeThing("Enter the new text...")

	// Call Service.
	obj, err := c.Service.Update(user, model.Text{Name: name, Tx: tx, Owner: user.GetModelUser().Email})

	if err != nil {
		c.DialogSrv.ShowErr(err)
		return
	}

	res, ok := obj.(*rpc.CreateResponse)
	if !ok {
		c.DialogSrv.ShowErr(errs.ErrTypeConversion)
		return
	}

	c.DialogSrv.ShowIt(fmt.Sprintf("%s %s\n\r", res.Status, res.UpdatedAt))
}

func (c *CommText) executeDelete(user user.User) {
	name := c.DialogSrv.GetDataAction("delete")

	// Call Service.
	obj, err := c.Service.Delete(user, model.Text{Name: name, Owner: user.GetModelUser().Email})

	if err != nil {
		c.DialogSrv.ShowErr(err)
		return
	}

	res, ok := obj.(*rpc.DeleteResponse)
	if !ok {
		c.DialogSrv.ShowErr(errs.ErrTypeConversion)
		return
	}
	c.DialogSrv.ShowIt(fmt.Sprintf("%s %s\n\r", res.Status, res.Name))
}
