package comm

import (
	"fmt"

	"github.com/boginskiy/GophKeeper/client/internal/errs"
	"github.com/boginskiy/GophKeeper/client/internal/infra"
	"github.com/boginskiy/GophKeeper/client/internal/model"
	"github.com/boginskiy/GophKeeper/client/internal/rpc"
	"github.com/boginskiy/GophKeeper/client/internal/service"
	"github.com/boginskiy/GophKeeper/client/internal/user"
)

type CommText struct {
	Dialoger infra.Dialoger
	Service  service.TextServicer[model.Text]
}

func NewCommText(
	dialoger infra.Dialoger,
	srv service.TextServicer[model.Text]) *CommText {

	return &CommText{
		Dialoger: dialoger,
		Service:  srv}
}

func (c *CommText) Registration(user user.User, dataType string) {
authLoop:
	for {
		// Get command.
		comm, _ := c.Dialoger.GetSomeThing(

			fmt.Sprintf("%s\n\r%s",
				"What do you want to do with the text: \n\r\t create \n\r\t read \n\r\t read-all \n\r\t update \n\r\t delete",
				"come back: back, need help: help"))

		switch comm {
		case "back", "help":
			break authLoop

		case "create":
			c.executeCreate(user, dataType)
		case "read":
			c.executeRead(user, dataType)
		case "read-all":
			c.executeReadAll(user, dataType)
		case "update":
			c.executeUpdate(user, dataType)
		case "delete":
			c.executeDelete(user, dataType)

		default:
			c.Dialoger.ShowIt("Invalid command. Try again...")
		}
	}
}

func (c *CommText) executeCreate(user user.User, dataType string) {
	name := c.Dialoger.GetDataAction("create")                 // Get info about name text.
	content, _ := c.Dialoger.GetSomeThing("Enter the text...") // Enter text for saving.

	// Call Service.
	obj, err := c.Service.Create(user, *model.NewText(name, dataType, content, user.GetModelUser().Email))

	if err != nil {
		c.Dialoger.ShowErr(err)
		return
	}

	res, ok := obj.(*rpc.CreateResponse)
	if !ok {
		c.Dialoger.ShowErr(errs.ErrTypeConversion)
		return
	}

	c.Dialoger.ShowIt(fmt.Sprintf("%s %s\n\r", res.Status, res.UpdatedAt))
}

func (c *CommText) executeRead(user user.User, dataType string) {
	name := c.Dialoger.GetDataAction("read")

	// Call Service.
	obj, err := c.Service.Read(user, model.Text{Name: name, Owner: user.GetModelUser().Email})

	if err != nil {
		c.Dialoger.ShowErr(err)
		return
	}

	res, ok := obj.(*rpc.ReadResponse)
	if !ok {
		c.Dialoger.ShowErr(errs.ErrTypeConversion)
		return
	}

	c.Dialoger.ShowIt(fmt.Sprintf("%s: %s   last update: %s\n\r", res.Name, res.Text, res.UpdatedAt))
}

func (c *CommText) executeReadAll(user user.User, dataType string) {
	// Call Service.
	obj, err := c.Service.ReadAll(user, model.Text{Type: dataType, Owner: user.GetModelUser().Email})

	if err != nil {
		c.Dialoger.ShowErr(err)
		return
	}

	res, ok := obj.(*rpc.ReadAllResponse)
	if !ok {
		c.Dialoger.ShowErr(errs.ErrTypeConversion)
		return
	}

	for _, text := range res.TextResponses {
		c.Dialoger.ShowIt(fmt.Sprintf("%s: %s   last update: %s", text.Name, text.Text, text.UpdatedAt))
	}
	fmt.Println()
}

func (c *CommText) executeUpdate(user user.User, dataType string) {
	name := c.Dialoger.GetDataAction("update")
	content, _ := c.Dialoger.GetSomeThing("Enter the new text...")

	// Call Service.
	obj, err := c.Service.Update(user, model.Text{Name: name, Content: content, Owner: user.GetModelUser().Email})

	if err != nil {
		c.Dialoger.ShowErr(err)
		return
	}

	res, ok := obj.(*rpc.CreateResponse)
	if !ok {
		c.Dialoger.ShowErr(errs.ErrTypeConversion)
		return
	}

	c.Dialoger.ShowIt(fmt.Sprintf("%s %s\n\r", res.Status, res.UpdatedAt))
}

func (c *CommText) executeDelete(user user.User, dataType string) {
	name := c.Dialoger.GetDataAction("delete")

	// Call Service.
	obj, err := c.Service.Delete(user, model.Text{Name: name, Owner: user.GetModelUser().Email})

	if err != nil {
		c.Dialoger.ShowErr(err)
		return
	}

	res, ok := obj.(*rpc.DeleteResponse)
	if !ok {
		c.Dialoger.ShowErr(errs.ErrTypeConversion)
		return
	}
	c.Dialoger.ShowIt(fmt.Sprintf("%s %s\n\r", res.Status, res.Name))
}
