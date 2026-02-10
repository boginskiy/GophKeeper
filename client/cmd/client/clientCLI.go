package client

import (
	"fmt"
	"os"
	"time"

	"github.com/boginskiy/GophKeeper/client/internal/logg"
)

const NAME = "CLIENT"

type ClientCLI struct {
	Logger logg.Logger
	Name   string
}

func NewClientCLI(logger logg.Logger) *ClientCLI {
	return &ClientCLI{
		Name:   NAME,
		Logger: logger,
	}
}

func (c *ClientCLI) print(text string) {
	for _, ch := range text {
		fmt.Fprintf(os.Stdout, "%c", ch)
		time.Sleep(15 * time.Millisecond)
	}
}

func (c *ClientCLI) SendMess(text ...string) {
	for _, t := range text {
		c.print(fmt.Sprintf("%s: %s\n\r", c.Name, t))
	}
}
